package config

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestProcessingInvalidBody(t *testing.T) {
	cfg := Load()
	log.Debugf("Config: %+v\n", cfg)
}
