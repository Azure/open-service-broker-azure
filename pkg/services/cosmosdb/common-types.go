package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

type provisioningParameters struct {
	IPFilterRules *ipFilterRule `json:"ipFilters"`
}

type ipFilterRule struct {
	Filters     []string `json:"allowedIPRanges"`
	AllowAzure  string   `json:"allowAccessFromAzure,omitempty"`
	AllowPortal string   `json:"allowAccessFromPortal,omitempty"`
}

type cosmosdbInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	DatabaseAccountName      string `json:"name"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	IPFilters                string `json:"ipFilters"`
}

type cosmosdbSecureInstanceDetails struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

// cosmosCredentials encapsulates CosmosDB-specific details for connecting via
// a variety of APIs. This excludes MongoDB.
type cosmosCredentials struct {
	URI                     string `json:"uri,omitempty"`
	PrimaryConnectionString string `json:"primaryConnectionString,omitempty"`
	PrimaryKey              string `json:"primaryKey,omitempty"`
}

func (c *cosmosAccountManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {

	pp := &provisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	return ppMap, nil, err
}

func (c *cosmosAccountManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}

func (
	c *cosmosAccountManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {
	p := map[string]service.ParameterSchema{}

	ipFilterSchema := make(map[string]service.ParameterSchema)

	allowAccessFromAzureSchema := service.NewParameterSchema(
		"string",
		"Specifies if Azure Services should be able to access"+
			" the CosmosDB account.",
	)
	allowAccessFromAzureSchema.SetAllowedValues(
		[]string{"", "enabled", "disabled"},
	)
	allowAccessFromAzureSchema.SetDefault("")
	ipFilterSchema["allowAccessFromAzure"] = allowAccessFromAzureSchema

	allowAccessFromPortalSchema := service.NewParameterSchema(
		"string",
		"Specifies if the Azure Portal should be able to"+
			" access the CosmosDB account. If `allowAccessFromAzure` is"+
			" set to enabled, this value is ignored.",
	)
	allowAccessFromPortalSchema.SetAllowedValues(
		[]string{"", "enabled", "disabled"},
	)
	allowAccessFromPortalSchema.SetDefault("")
	ipFilterSchema["allowAccessFromPortal"] = allowAccessFromPortalSchema

	allowedIPRangeSchema := service.NewParameterSchema(
		"array",
		"Values to include in IP Filter. Can be IP Address or CIDR range.",
	)
	err := allowedIPRangeSchema.SetItems(
		service.NewParameterSchema(
			"string",
			"Must be a valid IP address or CIDR",
		),
	)
	if err != nil {
		log.Errorf("error creating allowedIPRangeSchema: %s", err)
	}
	ipFilterSchema["allowedIPRanges"] = allowedIPRangeSchema

	ipFilters := service.NewParameterSchema(
		"object",
		"IP Range Filter to be applied to new CosmosDB account",
	)
	err = ipFilters.AddParameters(ipFilterSchema)
	if err != nil {
		log.Errorf("error adding properties to IP Filter schema: %s", err)
	}
	p["ipFilters"] = ipFilters

	return p

}
