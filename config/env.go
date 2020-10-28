package config

import (
	"github.com/kelseyhightower/envconfig"
)

type c struct {
	HTTP     http
	Service  service
	Redis    redis
	LogLevel string `envconfig:"LOGLEVEL"`
}

type redis struct {
	ClusterURL string `envconfig:"REDIS_CLUSTERURL"`
}

type http struct {
	Port string `envconfig:"HTTP_PORT"`
}

type service struct {
	Parser     string `envconfig:"SERVICE_PARSER"`
	City       string `envconfig:"SERVICE_CITY"`
	Category   string `envconfig:"SERVICE_CATEGORY"`
	Technology string `envconfig:"SERVICE_TECHNOLOGY"`
}

var Env c

func init() {
	envconfig.MustProcess("", &Env)
}
