package appinsights

import (
	"github.com/Azure/open-service-broker-azure/pkg/schemas"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func generateProvisioningParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{"location", "resourceGroup"},
		PropertySchemas: map[string]service.PropertySchema{
			"resourceGroup": schemas.GetResourceGroupSchema(),
			"location":      schemas.GetLocationSchema(),
			"appInsightsName": &service.StringPropertySchema{
				Title:       "Application Insights Name",
				Description: "The Application Insights component name",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}
