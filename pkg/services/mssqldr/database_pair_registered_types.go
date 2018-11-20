package mssqldr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *databasePairRegisteredManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return getDatabasePairRegisteredProvisionParamSchema()
}
