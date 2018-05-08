package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type provisioningParameters struct {
	IPFilterRules     *ipFilterRule      `json:"ipFilters"`
	ConsistencyPolicy *consistencyPolicy `json:"consistencyPolicy,omitempty"`
}

type ipFilterRule struct {
	Filters     []string `json:"allowedIPRanges"`
	AllowAzure  string   `json:"allowAccessFromAzure,omitempty"`
	AllowPortal string   `json:"allowAccessFromPortal,omitempty"`
}

type boundedStaleness struct {
	MaxStaleness *int `json:"maxStalenessPrefix,omitempty"`
	MaxInternal  *int `json:"maxIntervalInSeconds,omitempty"`
}

type consistencyPolicy struct {
	DefaultConsistency string            `json:"defaultConsistencyLevel,omitempty"`
	BoundedStaleness   *boundedStaleness `json:"boundedStaleness,omitempty"`
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
	URI                     string `json:"uri"`
	PrimaryConnectionString string `json:"primaryConnectionString"`
	PrimaryKey              string `json:"primaryKey"`
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
) getProvisionParametersSchema() service.InputParametersSchema {
	const maxStalenessPrefixMin = 1
	const maxStalenessPrefixMax = 2147483647
	const maxIntervalInSecondsMin = 5
	const maxIntervalInSecondsMax = 86400
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"ipFilters": &service.ObjectPropertySchema{
				Description: "IP Range Filter to be applied to new CosmosDB account",
				PropertySchemas: map[string]service.PropertySchema{
					"allowAccessFromAzure": &service.StringPropertySchema{
						Description: "Specifies if Azure Services should be able to access" +
							" the CosmosDB account.",
						AllowedValues: []string{"enabled", "disabled"},
						DefaultValue:  "enabled",
					},
					"allowAccessFromPortal": &service.StringPropertySchema{
						Description: "Specifies if the Azure Portal should be able to" +
							" access the CosmosDB account. If `allowAccessFromAzure` is" +
							" set to enabled, this value is ignored.",
						AllowedValues: []string{"enabled", "disabled"},
						DefaultValue:  "enabled",
					},
					"allowedIPRanges": &service.ArrayPropertySchema{
						Description: "Values to include in IP Filter. " +
							"Can be an IP Address or CIDR range.",
						ItemsSchema: &service.StringPropertySchema{
							Description: "Must be a valid IP address or CIDR",
						},
					},
				},
			},
			"consistencyPolicy": &service.ObjectPropertySchema{
				Description: "The consistency policy for the Cosmos DB account.",
				RequiredProperties: []string{
					"defaultConsistencyLevel",
				},
				PropertySchemas: map[string]service.PropertySchema{
					"defaultConsistencyLevel": &service.StringPropertySchema{
						Description: "The default consistency level and" +
							" configuration settings of the Cosmos DB account.",
						AllowedValues: []string{
							"Eventual",
							"Session",
							"BoundedStaleness",
							"Strong",
							"ConsistentPrefix",
						},
					},
					"boundedStaleness": &service.ObjectPropertySchema{
						Description: "The staleness settings when using " +
							"BoundedStaleness consistency.  Required when " +
							"using BoundedStaleness",
						PropertySchemas: map[string]service.PropertySchema{
							"maxStalenessPrefix": &service.IntPropertySchema{
								Description: "When used with the Bounded Staleness " +
									"consistency level, this value represents the number of " +
									"stale requests tolerated" +
									"Required when defaultConsistencyPolicy is set to " +
									" 'BoundedStaleness'.",
								MinValue: ptr.ToInt64(maxStalenessPrefixMin),
								MaxValue: ptr.ToInt64(maxStalenessPrefixMax),
							},
							"maxIntervalInSeconds": &service.IntPropertySchema{
								Description: "When used with the Bounded Staleness " +
									"consistency level, this value represents the time " +
									"amount of staleness (in seconds) tolerated. " +
									"Required when defaultConsistencyPolicy is set to " +
									" 'BoundedStaleness'.",
								MinValue: ptr.ToInt64(maxIntervalInSecondsMin),
								MaxValue: ptr.ToInt64(maxIntervalInSecondsMax),
							},
						},
					},
				},
			},
		},
	}
}
