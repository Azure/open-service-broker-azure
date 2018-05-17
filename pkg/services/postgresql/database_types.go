package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type databaseProvisioningParameters struct {
	Extensions []string `json:"extensions"`
}

func (
	d *databaseManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"extensions": dbExtensionsSchema,
		},
	}
}

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}
