package config

import (
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/kelseyhightower/envconfig"
)

// ModulesConfig represents details re: which modules should be included or
// excluded when the broker is started
type ModulesConfig interface {
	GetMinStability() service.Stability
}

type modulesConfig struct {
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"EXPERIMENTAL"`
	MinStability    service.Stability
}

// GetModulesConfig returns modules configuration
func GetModulesConfig() (ModulesConfig, error) {
	mc := modulesConfig{}
	err := envconfig.Process("", &mc)
	if err != nil {
		return mc, err
	}
	minStabilityStr := strings.ToUpper(mc.MinStabilityStr)
	switch minStabilityStr {
	case "EXPERIMENTAL":
		mc.MinStability = service.StabilityExperimental
	case "PREVIEW":
		mc.MinStability = service.StabilityPreview
	case "STABLE":
		mc.MinStability = service.StabilityStable
	default:
		return mc, fmt.Errorf(
			`unrecognized stability level "%s"`,
			minStabilityStr,
		)
	}
	return mc, nil
}

func (m modulesConfig) GetMinStability() service.Stability {
	return m.MinStability
}
