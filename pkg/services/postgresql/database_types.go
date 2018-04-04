package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type databaseProvisioningParameters struct {
	Extensions []string `json:"extensions"`
}

func (
	d *databaseManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {

	props := map[string]service.ParameterSchema{}
	props["extensions"] = &service.ArrayParameterSchema{
		Type:        "array",
		Description: "Database extensions to install",
		ItemsSchema: &service.SimpleParameterSchema{
			Type:        "string",
			Description: "Extension Name",
		},
	}
	return props
}

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

func (d *databaseManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := databaseProvisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}

	ppMap, err := service.GetMapFromStruct(pp)
	return ppMap, nil, err
}

func (d *databaseManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
