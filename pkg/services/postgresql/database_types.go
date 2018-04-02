package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

type databaseProvisioningParameters struct {
	Extensions []string `json:"extensions"`
}

func (
	d *databaseManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {

	props := map[string]service.ParameterSchema{}
	parentAliasSchema := service.NewParameterSchema(
		"string",
		"Specifies the alias of the DBMS upon which the database "+
			"should be provisioned.",
	)
	parentAliasSchema.SetRequired(true)
	props["parentAlias"] = parentAliasSchema

	extensionsSchema := service.NewParameterSchema(
		"array",
		"Database extensions to install",
	)
	err := extensionsSchema.SetItems(
		service.NewParameterSchema(
			"string",
			"Extension Name",
		),
	)
	if err != nil {
		log.Errorf("unable to build extensions schema: %s", err)
	}
	props["extensions"] = extensionsSchema
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
