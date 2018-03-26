package service

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// CatalogConfig represents details re: which modules' services should be
// included or excluded from the catalog
type CatalogConfig interface {
	GetMinStability() Stability
}

type catalogConfig struct {
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"EXPERIMENTAL"`
	MinStability    Stability
}

// GetCatalogConfigFromEnvironment returns catalog configuration
func GetCatalogConfigFromEnvironment() (CatalogConfig, error) {
	mc := catalogConfig{}
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

func (c catalogConfig) GetMinStability() Stability {
	return c.MinStability
}
