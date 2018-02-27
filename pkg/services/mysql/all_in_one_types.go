package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// AllInOneProvisioningParameters encapsulates non-sensitive dbms AND database
// MySQL-specific provisioning options
type AllInOneProvisioningParameters struct {
	DBMSProvisioningParameters `json:",squash"`
}

type allInOneInstanceDetails struct {
	dbmsInstanceDetails
	DatabaseName string `json:"database"`
}

type secureAllInOneInstanceDetails struct {
	secureDBMSInstanceDetails
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &AllInOneProvisioningParameters{}
}

func (
	a *allInOneManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return nil
}

func (
	a *allInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &allInOneInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &secureAllInOneInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return nil
}

func (
	a *allInOneManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return nil
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &bindingDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &secureBindingDetails{}
}
