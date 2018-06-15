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
	ARMDeploymentName          string               `json:"armDeployment"`
	FullyQualifiedDomainName   string               `json:"fullyQualifiedDomainName"` // nolint: lll
	ServerName                 string               `json:"server"`
	AdministratorLogin         string               `json:"administratorLogin"`
	AdministratorLoginPassword service.SecureString `json:"administratorLoginPassword"` // nolint: lll
}

func (d *dbmsManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbmsInstanceDetails{}
}

func (d *dbmsManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}
