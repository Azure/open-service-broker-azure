package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *dbmsManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return getDBMSCommonProvisionParamSchema()
}

func (
	d *dbmsManager,
) getUpdatingParametersSchema() service.InputParametersSchema {
	return getDBMSCommonUpdateParamSchema()
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	ServerName               string `json:"server"`
	AdministratorLogin       string `json:"administratorLogin"`
}

type secureDBMSInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}
