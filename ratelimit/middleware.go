package ratelimit

import (
	"github.com/nnqq/scr-api/logger"
	r "github.com/nnqq/scr-api/redis"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"net/http"
	"strings"
	"time"
)

var Middleware func(http.Handler) http.Handler

func init() {
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  10,
	}

	store, err := redis.NewStore(r.Client)
	logger.Must(err)

	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	bottleneck := stdlib.NewMiddleware(instance, stdlib.WithLimitReachedHandler(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte(`{"error":"Requests per second limit reached. Try again a bit later"}`))
			logger.Err(err)
		}),
	)

	Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Origin") == "https://leaq.ru" ||
				strings.HasPrefix(r.Header.Get("X-Real-Ip"), "10.") ||
				r.URL.Path == "/healthz" {
				// no rate limit for own SSR frontend, or k8s probe
				logger.Log.Debug().Str("path", r.URL.Path).Msg("no rate limit")
				next.ServeHTTP(w, r)
				return
			}

			logger.Log.Debug().Str("path", r.URL.Path).Msg("with rate limit")
			bottleneck.Handler(next).ServeHTTP(w, r)
		})
	}
}
