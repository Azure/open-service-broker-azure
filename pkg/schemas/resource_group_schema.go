package schemas

import "github.com/Azure/open-service-broker-azure/pkg/service"

// GetResourceGroupSchema returns pointer to general StringPropertySchema
// of "resourceGroup"
func GetResourceGroupSchema() *service.StringPropertySchema {
	return &service.StringPropertySchema{
		Title: "Resource group",
		Description: "The (new or existing) resource group with which" +
			" to associate new resources.",
	}
}
