// +build experimental

package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type allInOneInstanceDetails struct {
	dbmsInstanceDetails `json:",squash"`
	DatabaseName        string `json:"database"`
}

type secureAllInOneInstanceDetails struct {
	secureDBMSInstanceDetails `json:",squash"`
}

func (a *allInOneManager) SplitProvisioningParameters(
	cpp map[string]interface{},
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	// nolint: megacheck
	pp := struct {
		dbmsProvisioningParams  `json:",squash"`
		databaseProvisionParams `json:",squash"`
	}{}

	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	return ppMap, nil, err
}

func (a *allInOneManager) SplitBindingParameters(
	params map[string]interface{},
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
