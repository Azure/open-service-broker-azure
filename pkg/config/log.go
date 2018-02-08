package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// LogConfig represents configuration options for the broker's leveled logging
type LogConfig interface {
	GetLevel() log.Level
}

type logConfig struct {
	LevelStr string `envconfig:"LOG_LEVEL" default:"INFO"`
	Level    log.Level
}

// GetLogConfig returns log configuration
func GetLogConfig() (LogConfig, error) {
	lc := logConfig{}
	err := envconfig.Process("", &lc)
	if err != nil {
		return lc, err
	}
	lc.Level, err = log.ParseLevel(lc.LevelStr)
	return lc, err
}

func (l logConfig) GetLevel() log.Level {
	return l.Level
}
