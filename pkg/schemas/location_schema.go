package schemas

import "github.com/Azure/open-service-broker-azure/pkg/service"

// GetLocationSchema returns pointer to general StringPropertySchema
// of "location"
func GetLocationSchema() *service.StringPropertySchema {
	return &service.StringPropertySchema{
		Title: "Location",
		Description: "The Azure region in which to provision" +
			" applicable resources.",
		OneOf: []service.EnumValue{
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
			{Value: "japaneast", Title: "Japan East"},
			{Value: "japanwest", Title: "Japan West"},
			{Value: "koreacentral", Title: "Korea Central"},
			{Value: "koreasouth", Title: "Korea South"},
			{Value: "northcentralus", Title: "North Central US"},
			{Value: "northeurope", Title: "North Europe"},
			{Value: "southcentralus", Title: "South Central US"},
			{Value: "southeastasia", Title: "Southeast Asia"},
			{Value: "southindia", Title: "South India"},
			{Value: "uksouth", Title: "UK South"},
			{Value: "ukwest", Title: "UK West"},
			{Value: "westcentralus", Title: "West Central US"},
			{Value: "westeurope", Title: "West Europe"},
			{Value: "westindia", Title: "West India"},
			{Value: "westus", Title: "West US"},
			{Value: "westus2", Title: "West US 2"},
		},
	}
}
