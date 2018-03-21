package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databaseProvisioningParameters struct {
	Extensions []string `json:"extensions"`
}

// GetDatabaseProvisionParametersSchema generates a schema for parameters used
// in database instance provisioning
func GetDatabaseProvisionParametersSchema() *service.ParametersSchema {
	p := service.GetEmptyParameterSchema()
	props := map[string]interface{}{}
	props["parentAlias"] = service.Parameter{
		Type: "string",
		Description: "Specifies the alias of the DBMS upon which the database " +
			"should be provisioned.",
	}
	props["extensions"] = service.Parameter{
		Type: "array",
		Items: service.Parameter{
			Type: "string",
		},
	}
	p.Properties = props
	p.Required = []string{"parentAlias"}
	return p
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
