package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port string `env:"PORT"`
}

func Load() (cfg Config) {
	if err := env.Parse(&cfg); err != nil {
		log.Errorf("%s", err)
	}
	return
}
