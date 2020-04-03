package schemas

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

var (
	azurePublicCloudLocationSchema = service.StringPropertySchema{
		Title: "Location",
		Description: "The Azure region in which to provision" +
			" applicable resources.",
		OneOf: []service.EnumValue{
			{Value: "australiacentral", Title: "Australia Central"},
			{Value: "australiacentral2", Title: "Australia Central 2"},
			{Value: "australiaeast", Title: "Australia East"},
			{Value: "australiasoutheast", Title: "Australia Southeast"},
			{Value: "brazilsouth", Title: "Brazil South"},
			{Value: "canadacentral", Title: "Canada Central"},
			{Value: "canadaeast", Title: "Canada East"},
			{Value: "centralindia", Title: "Central India"},
			{Value: "centralus", Title: "Central US"},
			{Value: "eastasia", Title: "East Asia"},
			{Value: "eastus", Title: "East US"},
			{Value: "eastus2", Title: "East US 2"},
			{Value: "francecentral", Title: "France Central"},
			{Value: "francesouth", Title: "France South"},
			{Value: "germanynorth", Title: "Germany North"},
			{Value: "germanywestcentral", Title: "Germany West Central"},
			{Value: "japaneast", Title: "Japan East"},
			{Value: "japanwest", Title: "Japan West"},
			{Value: "koreacentral", Title: "Korea Central"},
			{Value: "koreasouth", Title: "Korea South"},
			{Value: "northcentralus", Title: "North Central US"},
			{Value: "northeurope", Title: "North Europe"},
			{Value: "norwayeast", Title: "Norway East"},
			{Value: "norwaywest", Title: "Norway West"},
			{Value: "southafricanorth", Title: "South Africa North"},
			{Value: "southafricawest", Title: "South Africa West"},
			{Value: "southcentralus", Title: "South Central US"},
			{Value: "southeastasia", Title: "Southeast Asia"},
			{Value: "southindia", Title: "South India"},
			{Value: "switzerlandnorth", Title: "Switzerland North"},
			{Value: "switzerlandwest", Title: "Switzerland West"},
			{Value: "uaecentral", Title: "UAE Central"},
			{Value: "uaenorth", Title: "UAE North"},
			{Value: "uksouth", Title: "UK South"},
			{Value: "ukwest", Title: "UK West"},
			{Value: "westcentralus", Title: "West Central US"},
			{Value: "westeurope", Title: "West Europe"},
			{Value: "westindia", Title: "West India"},
			{Value: "westus", Title: "West US"},
			{Value: "westus2", Title: "West US 2"},
		},
	}
	azureChinaCloudLocationSchema = service.StringPropertySchema{
		Title: "Location",
		Description: "The Azure region in which to provision" +
			" applicable resources.",
		OneOf: []service.EnumValue{
			{Value: "chinanorth2", Title: "China North 2"},
			{Value: "chinanorth", Title: "China North"},
			{Value: "chinaeast2", Title: "China East 2"},
			{Value: "chinaeast", Title: "China East"},
		},
	}
)

// GetLocationSchema returns pointer to general StringPropertySchema
// of "location"
func GetLocationSchema() *service.StringPropertySchema {
	environmentName := azure.GetEnvironmentName()
	switch environmentName {
	case "AzureChinaCloud":
		return &azureChinaCloudLocationSchema
	case "AzurePublicCloud":
		return &azurePublicCloudLocationSchema
	}
	// We shouldn't run into here, this function is supposed
	// to return at the `switch` statement. We return
	// public location schema here so that the error can be
	// reported when provisioning the resource.
	return &azurePublicCloudLocationSchema
}
