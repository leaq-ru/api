package config

import (
	"github.com/kelseyhightower/envconfig"
)

type c struct {
	HTTP     http
	Service  service
	LogLevel string `envconfig:"LOGLEVEL"`
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
