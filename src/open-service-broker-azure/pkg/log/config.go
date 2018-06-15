package log

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// Config represents configuration options for the broker's leveled logging
type Config interface {
	GetLevel() log.Level
}

type config struct {
	LevelStr string `envconfig:"LOG_LEVEL" default:"INFO"`
	Level    log.Level
}

// GetConfig returns log configuration
func GetConfig() (Config, error) {
	lc := config{}
	err := envconfig.Process("", &lc)
	if err != nil {
		return lc, err
	}
	lc.Level, err = log.ParseLevel(lc.LevelStr)
	return lc, err
}

func (c config) GetLevel() log.Level {
	return c.Level
}
