package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

func (
	d *databaseManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {
	props := map[string]service.ParameterSchema{}
	parentAliasSchema := service.NewSimpleParameterSchema(
		"string",
		"Specifies the alias of the DBMS upon which the database "+
			"should be provisioned.",
	)
	parentAliasSchema.SetRequired(true)
	props["parentAlias"] = parentAliasSchema
	return props
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
