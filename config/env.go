package config

import "github.com/kelseyhightower/envconfig"

type c struct {
	HTTP             http
	Service          service
	Redis            redis
	Robokassa        robokassa
	DisableRateLimit string `envconfig:"DISABLERATELIMIT"`
	LogLevel         string `envconfig:"LOGLEVEL"`
}

type robokassa struct {
	WebhookSecret string `envconfig:"ROBOKASSA_WEBHOOKSECRET"`
}

type redis struct {
	URL string `envconfig:"REDIS_URL"`
}

type http struct {
	Port string `envconfig:"HTTP_PORT"`
}

type service struct {
	Parser     string `envconfig:"SERVICE_PARSER"`
	City       string `envconfig:"SERVICE_CITY"`
	Category   string `envconfig:"SERVICE_CATEGORY"`
	Technology string `envconfig:"SERVICE_TECHNOLOGY"`
	User       string `envconfig:"SERVICE_USER"`
	Billing    string `envconfig:"SERVICE_BILLING"`
	Exporter   string `envconfig:"SERVICE_EXPORTER"`
	Org        string `envconfig:"SERVICE_ORG"`
}

var Env c

func init() {
	envconfig.MustProcess("", &Env)
}
