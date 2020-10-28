package redis

import (
	rd "github.com/go-redis/redis/v7"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
)

var Client *rd.ClusterClient

func init() {
	rdb := rd.NewClusterClient(&rd.ClusterOptions{
		Addrs: []string{config.Env.Redis.ClusterURL},
	})

	err := rdb.Ping().Err()
	logger.Must(err)

	Client = rdb
}
