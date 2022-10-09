package ratelimit

import (
	"context"
	"fmt"
	"github.com/leaq-ru/api/call"
	"github.com/leaq-ru/api/config"
	"github.com/leaq-ru/api/logger"
	"github.com/leaq-ru/api/middleware"
	r "github.com/leaq-ru/api/redis"
	"github.com/leaq-ru/proto/codegen/go/billing"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"net/http"
	"strings"
	"time"
)

var Middleware func(http.Handler) http.Handler

func init() {
	store, e := redis.NewStore(r.Client)
	logger.Must(e)

	Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const headerDataPremium = "Grpc-Metadata-Data-Premium"
			r.Header.Set(headerDataPremium, "")

			userID := r.Header.Get(middleware.HeaderUserID)
			var premium bool
			if userID != "" {
				ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
				defer cancel()

				plan, err := call.Billing.GetDataPlan(ctx, &billing.GetDataPlanRequest{
					UserId: userID,
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

				if plan.GetPremium() {
					premium = plan.GetPremium()
					r.Header.Set(headerDataPremium, "true")
				}
			}

			origin := r.Header.Get("Origin")
			xRealIp := r.Header.Get("X-Real-Ip")
			xForwardedFor := r.Header.Get("X-Forwarded-For")
			path := r.URL.Path

			dbg := logger.Log.Debug().
				Str("origin", origin).
				Str("xRealIp", xRealIp).
				Str("xForwardedFor", xForwardedFor).
				Str("path", path).
				Bool("premium", premium)

			if config.Env.DisableRateLimit == "true" ||
				origin == "https://leaq.ru" ||
				xRealIp == "" ||
				strings.HasPrefix(path, "/docs/") ||
				path == "/v1/billing/robokassaWebhook/"+config.Env.Robokassa.WebhookSecret {
				// no rate limit for own frontend, or k8s probe, or Robokassa webhook with valid secret
				dbg.Msg("no rate limit")
				next.ServeHTTP(w, r)
				return
			}

			dbg.Msg("with rate limit")

			var rateRPS limiter.Rate
			if premium {
				rateRPS = limiter.Rate{
					Limit:  30,
					Period: time.Second,
				}
			} else {
				rateRPS = limiter.Rate{
					Limit:  10,
					Period: time.Second,
				}
			}

			opts := limiter.WithTrustForwardHeader(true)

			rps := stdlib.NewMiddleware(limiter.New(store, rateRPS, opts), makeLimitReached("second"))
			rps.Handler(next).ServeHTTP(w, r)
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
