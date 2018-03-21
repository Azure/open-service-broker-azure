package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

// GetDatabaseProvisionParametersSchema generates a schema for parameters used
// in database instance provisioning
func GetDatabaseProvisionParametersSchema() *service.ParametersSchema {
	p := service.GetEmptyParameterSchema()
	props := map[string]service.Parameter{}
	props["parentAlias"] = service.Parameter{
		Type: "string",
		Description: "Specifies the alias of the DBMS upon which the database " +
			"should be provisioned.",
	}
	p.Properties = props
	p.Required = []string{"parentAlias"}
	return p
}

func (d *databaseManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	return nil, nil, nil
}

func (d *databaseManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
