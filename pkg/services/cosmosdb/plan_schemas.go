package cosmosdb

import (
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func generateProvisioningParamsSchema() service.InputParametersSchema {
	const maxStalenessPrefixMin = 1
	const maxStalenessPrefixMax = 2147483647
	const maxIntervalInSecondsMin = 5
	const maxIntervalInSecondsMax = 86400
	return service.InputParametersSchema{
		RequiredProperties: []string{"location", "resourceGroup"},
		PropertySchemas: map[string]service.PropertySchema{
			"location": &service.StringPropertySchema{
				Title: "Location",
				Description: "The Azure region in which to provision" +
					" applicable resources.",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"resourceGroup": &service.StringPropertySchema{
				Title: "Resource group",
				Description: "The (new or existing) resource group with which" +
					" to associate new resources.",
			},
			"ipFilters": &service.ObjectPropertySchema{
				Title:       "IP filters",
				Description: "IP Range Filter to be applied to new CosmosDB account",
				PropertySchemas: map[string]service.PropertySchema{
					"allowAccessFromAzure": &service.StringPropertySchema{
						Title: "Allow access from Azure",
						Description: "Specifies if Azure Services should be able to access" +
							" the CosmosDB account.",
						AllowedValues: []string{"enabled", "disabled"},
						DefaultValue:  "enabled",
					},
					"allowAccessFromPortal": &service.StringPropertySchema{
						Title: "Allow access From Portal",
						Description: "Specifies if the Azure Portal should be able to" +
							" access the CosmosDB account. If `allowAccessFromAzure` is" +
							" set to enabled, this value is ignored.",
						AllowedValues: []string{"enabled", "disabled"},
						DefaultValue:  "enabled",
					},
					"allowedIPRanges": &service.ArrayPropertySchema{
						Title: "Allowed IP ranges",
						Description: "Values to include in IP Filter. " +
							"Can be an IP Address or CIDR range.",
						ItemsSchema: &service.StringPropertySchema{
							Description:             "Must be a valid IP address or CIDR",
							CustomPropertyValidator: ipRangeValidator,
						},
					},
				},
				DefaultValue: map[string]interface{}{
					"allowAccessFromAzure": "enabled",
				},
			},
			"consistencyPolicy": &service.ObjectPropertySchema{
				Title:       "Consistency policy",
				Description: "The consistency policy for the Cosmos DB account.",
				RequiredProperties: []string{
					"defaultConsistencyLevel",
				},
				PropertySchemas: map[string]service.PropertySchema{
					"defaultConsistencyLevel": &service.StringPropertySchema{
						Title: "Default consistency level",
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
						Title: "Bounded staleness",
						Description: "The staleness settings when using " +
							"BoundedStaleness consistency.  Required when " +
							"using BoundedStaleness",
						RequiredProperties: []string{
							"maxStalenessPrefix",
							"maxIntervalInSeconds",
						},
						PropertySchemas: map[string]service.PropertySchema{
							"maxStalenessPrefix": &service.IntPropertySchema{
								Title: "Maximum staleness prefix",
								Description: "When used with the Bounded Staleness " +
									"consistency level, this value represents the number of " +
									"stale requests tolerated" +
									"Required when defaultConsistencyPolicy is set to " +
									" 'BoundedStaleness'.",
								MinValue: ptr.ToInt64(maxStalenessPrefixMin),
								MaxValue: ptr.ToInt64(maxStalenessPrefixMax),
							},
							"maxIntervalInSeconds": &service.IntPropertySchema{
								Title: "Maximum interval in seconds",
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
				CustomPropertyValidator: consistencyPolicyValidator,
			},
		},
	}
}

func ipRangeValidator(context, value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		cidr, _, _ := net.ParseCIDR(value)
		if cidr == nil {
			return service.NewValidationError(
				context,
				fmt.Sprintf(`"%s" is neither a valid IP address or CIDR range`, value),
			)
		}
	}
	return nil
}

// consistencyPolicyValidator is used to validate that the boundedStaleness
// object has been included in the parameters if "BoundedStaleness" was the
// selected default consistency policy. No further validation of the
// boundedStaleness (other than checking that it exists) is carried out here
// because the schema-based validations have all the rest of that covered.
func consistencyPolicyValidator(
	context string,
	valMap map[string]interface{},
) error {
	defaultConsistencyLevel := valMap["defaultConsistencyLevel"].(string)
	if defaultConsistencyLevel == "BoundedStaleness" {
		_, ok := valMap["boundedStaleness"].(map[string]interface{})
		if !ok {
			return service.NewValidationError(
				fmt.Sprintf("%s.boundedStaleness", context),
				"field is required",
			)
		}
	}
	return nil
}
