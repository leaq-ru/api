package redis

import (
	rd "github.com/go-redis/redis/v7"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
)

var Client *rd.Client

func init() {
	Client = rd.NewClient(&rd.Options{
		Addr: config.Env.Redis.URL,
	})

	err := Client.Ping().Err()
	logger.Must(err)
}
