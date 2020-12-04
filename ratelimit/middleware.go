package ratelimit

import (
	"context"
	"fmt"
	"github.com/nnqq/scr-api/call"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-api/middleware"
	r "github.com/nnqq/scr-api/redis"
	"github.com/nnqq/scr-proto/codegen/go/billing"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"net/http"
	"strings"
	"time"
)

var Middleware func(http.Handler) http.Handler

func init() {
	store, err := redis.NewStoreWithOptions(r.Client, limiter.StoreOptions{
		Prefix: "api",
	})
	logger.Must(err)

	Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const headerDataPremium = "Grpc-Metadata-data-premium"
			r.Header.Set(headerDataPremium, "")

			origin := r.Header.Get("Origin")
			path := r.URL.Path

			if origin == "https://leaq.ru" ||
				strings.HasPrefix(r.Header.Get("X-Real-Ip"), "10.") ||
				strings.HasPrefix(path, "/docs/") ||
				path == "/healthz" {
				// no rate limit for own frontend, or k8s probe
				logger.Log.Debug().Str("path", path).Msg("no rate limit")
				next.ServeHTTP(w, r)
				return
			}

			logger.Log.Debug().Str("path", path).Msg("with rate limit")

			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			plan, err := call.Billing.GetDataPlan(ctx, &billing.GetDataPlanRequest{
				UserId: r.Header.Get(middleware.HeaderUserID),
			})
			if err != nil {
				logger.Log.Error().Err(err).Send()
				w.WriteHeader(http.StatusUnauthorized)

				_, err = w.Write(nil)
				if err != nil {
					logger.Log.Error().Err(err).Send()
				}
				return
			}

			if plan.Premium {
				r.Header.Set(headerDataPremium, "true")
			}

			var rateRPS limiter.Rate
			if plan.Premium {
				rateRPS = limiter.Rate{
					Limit:  30,
					Period: time.Second,
				}
			} else {
				rateRPS = limiter.Rate{
					Limit:  3,
					Period: time.Second,
				}
			}

			opts := limiter.WithTrustForwardHeader(true)

			bottleneckRPS := stdlib.NewMiddleware(limiter.New(store, rateRPS, opts), makeLimitReached("second"))
			if plan.Premium {
				bottleneckRPS.Handler(next).ServeHTTP(w, r)
				return
			}

			rateRPD := limiter.Rate{
				Limit:  1000,
				Period: 24 * time.Hour,
			}
			bottleneckRPD := stdlib.NewMiddleware(limiter.New(store, rateRPD, opts), makeLimitReached("day"))
			bottleneckRPD.Handler(bottleneckRPS.Handler(next)).ServeHTTP(w, r)
		})
	}
}

func makeLimitReached(interval string) stdlib.Option {
	return stdlib.WithLimitReachedHandler(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, err := w.Write([]byte(fmt.Sprintf(
			`{"error":"Requests per %s limit reached. Try again a bit later"}`,
			interval,
		)))
		logger.Err(err)
	})
}
