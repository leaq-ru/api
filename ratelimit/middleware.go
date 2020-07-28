package ratelimit

import (
	"github.com/nnqq/scr-api/logger"
	r "github.com/nnqq/scr-api/redis"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"net/http"
	"time"
)

var Middleware *stdlib.Middleware

func init() {
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  30,
	}

	store, err := redis.NewStore(r.Client)
	logger.Must(err)

	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	Middleware = stdlib.NewMiddleware(instance, stdlib.WithLimitReachedHandler(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, err := w.Write([]byte(`{"error":"Requests per second limit reached. Try again a bit later"}`))
		logger.Log.Error().Err(err).Send()
	}))
}
