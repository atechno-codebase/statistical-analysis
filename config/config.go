package config

import (
	"github.com/AbsaOSS/env-binder/env"
)

type Config struct {
	Port     string `env:"PORT,default=8000"`
	MongoURL string `env:"MONGO_URL"`
	Database string `env:"DATABASE"`
}

var Configuration Config

func Init() error {
	return env.Bind(&Configuration)
}
