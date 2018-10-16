package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// CatalogConfig represents details re: which modules' services should be
// included or excluded from the catalog
type CatalogConfig struct {
	MinStability            Stability
	UseV2Guid               bool
	EnableMigrationServices bool
}

type tempCatalogConfig struct {
	CatalogConfig
	MinStabilityStr            string `envconfig:"MIN_STABILITY" default:"PREVIEW"`
	UseV2GuidStr               string `envconfig:"USE_V2_GUID" default:"false"`
	EnableMigrationServicesStr string `envconfig:"ENABLE_MIGRATION_SERVICES" default:"false"` // nolint: lll
}

// NewCatalogConfigWithDefaults returns a CatalogConfig object with default
// values already applied. Callers are then free to set custom values for the
// remaining fields and/or override default values.
func NewCatalogConfigWithDefaults() CatalogConfig {
	return CatalogConfig{
		MinStability:            StabilityPreview,
		UseV2Guid:               false,
		EnableMigrationServices: false,
	}
}

// GetCatalogConfigFromEnvironment returns catalog configuration
func GetCatalogConfigFromEnvironment() (CatalogConfig, error) {
	c := tempCatalogConfig{
		CatalogConfig: NewCatalogConfigWithDefaults(),
	}
	err := envconfig.Process("", &c)
	if err != nil {
		return c.CatalogConfig, err
	}
	minStabilityStr := strings.ToUpper(c.MinStabilityStr)
	switch minStabilityStr {
	case "EXPERIMENTAL":
		c.MinStability = StabilityExperimental
	case "PREVIEW":
		c.MinStability = StabilityPreview
	case "STABLE":
		c.MinStability = StabilityStable
	default:
		return c.CatalogConfig, fmt.Errorf(
			`unrecognized stability level "%s"`,
			minStabilityStr,
		)
	}
	c.UseV2Guid, err = strconv.ParseBool(c.UseV2GuidStr)
	if err != nil {
		return c.CatalogConfig, fmt.Errorf(
			`unrecognized UseV2Guid boolean "%s": %s`,
			c.EnableMigrationServicesStr,
			err,
		)
	}
	c.EnableMigrationServices, err =
		strconv.ParseBool(c.EnableMigrationServicesStr)
	if err != nil {
		return c.CatalogConfig, fmt.Errorf(
			`unrecognized EnableMigrationServices boolean "%s": %s`,
			c.EnableMigrationServicesStr,
			err,
		)
	}
	return c.CatalogConfig, nil
}
