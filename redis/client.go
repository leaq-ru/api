package redis

import (
	rd "github.com/go-redis/redis/v7"
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
	"time"
)

var Client *rd.ClusterClient

func init() {
	rdb := rd.NewClusterClient(&rd.ClusterOptions{
		Addrs: []string{config.Env.Redis.ClusterURL},
	})

	err := rdb.Ping().Err()
	logger.Must(err)

	Client = rdb

	go func() {
		for {
			time.Sleep(time.Second)

			e := Client.Ping().Err()
			if e != nil {
				logger.Log.Error().Err(e).Send()
				Client = rd.NewClusterClient(&rd.ClusterOptions{
					Addrs: []string{config.Env.Redis.ClusterURL},
				})
			}
		}
	}()
}
