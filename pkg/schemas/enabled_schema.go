package schemas

import "github.com/Azure/open-service-broker-azure/pkg/service"

const (
	// EnabledParamString defines "enabled" value
	EnabledParamString = "enabled"
	// DisabledParamString defines "enabled" value
	DisabledParamString = "disabled"
)

// EnabledDisabledValues returns "enabled" and "disabled" EnumValues
func EnabledDisabledValues() []service.EnumValue {
	return []service.EnumValue{
		{
			Value: EnabledParamString,
			Title: "Enabled",
		},
		{
			Value: DisabledParamString,
			Title: "Disabled",
		},
	}
}
