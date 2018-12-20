package schemas

import "github.com/Azure/open-service-broker-azure/pkg/service"

// GetLocationSchema returns pointer to general StringPropertySchema
// of "location"
func GetLocationSchema() *service.StringPropertySchema {
	return &service.StringPropertySchema{
		Title: "Location",
		Description: "The Azure region in which to provision" +
			" applicable resources.",
		AllowedValues: []string{
			"australiaeast",
			"australiasoutheast",
			"brazilsouth",
			"canadacentral",
			"canadaeast",
			"centralindia",
			"centralus",
			"eastasia",
			"eastus",
			"eastus2",
			"japaneast",
			"japanwest",
			"koreacentral",
			"koreasouth",
			"northcentralus",
			"northeurope",
			"southcentralus",
			"southeastasia",
			"southindia",
			"uksouth",
			"ukwest",
			"westcentralus",
			"westeurope",
			"westindia",
			"westus",
			"westus2",
		},
	}
}
