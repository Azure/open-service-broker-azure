package mssqldr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *dbmsPairRegisteredManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return getDBMSPairRegisteredProvisionParamSchema()
}

func (
	d *dbmsPairRegisteredManager,
) getUpdatingParametersSchema() service.InputParametersSchema {
	return getDBMSPairRegisteredUpdateParamSchema()
}

func (d *dbmsPairRegisteredManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &dbmsPairInstanceDetails{}
}

func (d *dbmsPairRegisteredManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return nil
}

// nolint: lll
type dbmsPairInstanceDetails struct {
	PriARMDeploymentName          string               `json:"primaryArmDeployment"`
	PriFullyQualifiedDomainName   string               `json:"primaryFullyQualifiedDomainName"`
	PriServerName                 string               `json:"primaryServer"`
	PriAdministratorLogin         string               `json:"primaryAdministratorLogin"`
	PriAdministratorLoginPassword service.SecureString `json:"primaryAdministratorLoginPassword"`
	SecARMDeploymentName          string               `json:"secondaryArmDeployment"`
	SecFullyQualifiedDomainName   string               `json:"secondaryFullyQualifiedDomainName"`
	SecServerName                 string               `json:"secondaryServer"`
	SecAdministratorLogin         string               `json:"secondaryAdministratorLogin"`
	SecAdministratorLoginPassword service.SecureString `json:"secondaryAdministratorLoginPassword"`
}
