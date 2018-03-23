package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	IPFilterRules *ipFilterRule `json:"ipFilters"`
}

type ipFilterRule struct {
	Filters     []string `json:"filters"`
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
	return nil, nil, nil
}

func (c *cosmosAccountManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
