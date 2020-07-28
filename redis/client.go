package redis

import (
	r "github.com/go-redis/redis/v7"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
)

var Client *r.Client

func init() {
	rdb := r.NewClient(&r.Options{
		Addr: config.Env.Redis.URL,
	})

	err := rdb.Ping().Err()
	logger.Must(err)

	Client = rdb
}
