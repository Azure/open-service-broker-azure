package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type dbmsInstanceDetails struct {
	ARMDeploymentName          string               `json:"armDeployment"`
	ServerName                 string               `json:"server"`
	FullyQualifiedDomainName   string               `json:"fullyQualifiedDomainName"`   // nolint: lll
	AdministratorLogin         string               `json:"administratorLogin"`         // nolint: lll
	AdministratorLoginPassword service.SecureString `json:"administratorLoginPassword"` // nolint: lll
}

func (d *dbmsManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbmsInstanceDetails{}
}

func (d *dbmsManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}
