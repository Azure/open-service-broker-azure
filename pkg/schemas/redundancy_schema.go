package schemas

import "github.com/Azure/open-service-broker-azure/pkg/service"

// LocalRedundancy provides an EnumValue slice with "local" value.
func LocalRedundancy() []service.EnumValue {
	return []service.EnumValue{{
		Value: "local",
		Title: "Local",
	}}
}

// LocalAndGeoRedundancy provides an EnumValue slice with
// "local" and "geo" value.
func LocalAndGeoRedundancy() []service.EnumValue {
	return []service.EnumValue{
		{Value: "local", Title: "Local"},
		{Value: "geo", Title: "Geo"},
	}
}
