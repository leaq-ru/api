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
	URL    string `envconfig:"REDIS_URL"`
	Master string `envconfig:"REDIS_MASTER"`
}

type http struct {
	Port string `envconfig:"HTTP_PORT"`
}

type service struct {
	Parser   string `envconfig:"SERVICE_PARSER"`
	City     string `envconfig:"SERVICE_CITY"`
	Category string `envconfig:"SERVICE_CATEGORY"`
}

var Env c

func init() {
	envconfig.MustProcess("", &Env)
}
