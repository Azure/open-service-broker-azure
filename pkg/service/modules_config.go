package service

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// ModulesConfig represents details re: which modules' services should be
// included or excluded from the catalog
type ModulesConfig interface {
	GetMinStability() Stability
}

type modulesConfig struct {
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"EXPERIMENTAL"`
	MinStability    Stability
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
		mc.MinStability = StabilityExperimental
	case "PREVIEW":
		mc.MinStability = StabilityPreview
	case "STABLE":
		mc.MinStability = StabilityStable
	default:
		return mc, fmt.Errorf(
			`unrecognized stability level "%s"`,
			minStabilityStr,
		)
	}
	return mc, nil
}

func (m modulesConfig) GetMinStability() Stability {
	return m.MinStability
}
