package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

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
) getProvisionParametersSchema() map[string]*service.ParameterSchema {
	p := map[string]*service.ParameterSchema{}

	ipFilterSchema := make(map[string]*service.ParameterSchema)
	ipFilterSchema["allowAccessFromAzure"] = &service.ParameterSchema{
		Type: "string",
		Description: "Specifies if Azure Services should be able to access" +
			" the CosmosDB account. Valid valued are ``, `enabled`, or " +
			"`disabled`.",
	}

	ipFilterSchema["allowAccessFromPortal"] = &service.ParameterSchema{
		Type: "string",
		Description: "Specifies if the Azure Portal should be able to" +
			" access the CosmosDB account. If `allowAccessFromAzure` is" +
			" set to enabled, this value is ignored. Valid valued are" +
			" `` (unspecified), `enabled`, or `disabled`.",
	}

	ipFilterSchema["allowedIPRanges"] = &service.ParameterSchema{
		Type: "array",
		Description: "Values to include in IP Filter. Can be IP Address or" +
			"CIDR range.",
		Items: &service.ParameterSchema{
			Type:        "string",
			Description: "Must be a valid IP address or CIDR",
		},
	}

	p["ipFilters"] = &service.ParameterSchema{
		Type:        "object",
		Description: "IP Range Filter to be applied to new CosmosDB account",
		Properties:  ipFilterSchema,
	}

	return p

}
