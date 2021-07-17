package redis

import (
	rd "github.com/go-redis/redis/v7"
	"github.com/leaq-ru/api/config"
	"github.com/leaq-ru/api/logger"
)

var Client *rd.Client

func init() {
	Client = rd.NewClient(&rd.Options{
		Addr: config.Env.Redis.URL,
	})

	err := Client.Ping().Err()
	logger.Must(err)
}
