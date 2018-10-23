// +build !unit

package lifecycle

import (
	"context"
	"fmt"

	networkSDK "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network" // nolint: lll
	storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage" // nolint: lll
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"

	uuid "github.com/satori/go.uuid"
)

var rediscacheTestCases = []serviceLifecycleTestCase{
	{
		group:     "rediscache",
		name:      "rediscache-basic-provision-and-update",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
		provisioningParameters: map[string]interface{}{
			"location":         "southcentralus",
			"skuCapacity":      1,
			"enableNonSslPort": "disabled",
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
		updatingParameters: map[string]interface{}{
			"skuCapacity":      2,
			"enableNonSslPort": "enabled",
			"tags": map[string]interface{}{
				"latest-operation": "update",
			},
		},
	},
	{
		group:     "rediscache",
		name:      "rediscache-premium-shard-count-provision-and-update",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
		provisioningParameters: map[string]interface{}{
			"location":   "eastus",
			"shardCount": 2,
		},
		updatingParameters: map[string]interface{}{
			"shardCount": 1,
		},
	},
	{
		group:     "rediscache",
		name:      "rediscache-premium-provision-and-update",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		updatingParameters: map[string]interface{}{
			"skuCapacity":      1,
			"enableNonSslPort": "enabled",
		},
		preProvisionFns: []preProvisionFn{
			createVirtualNetwork,
		},
	},
	{
		group:     "rediscache",
		name:      "rediscache-premium-rdb-backup-test",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
		updatingParameters: map[string]interface{}{
			"enableNonSslPort": "disabled",
			"redisConfiguration": map[string]interface{}{
				"rdb-backup-enabled": "disabled",
			},
		},
		preProvisionFns: []preProvisionFn{
			createStorageAccount,
		},
	},
}

// nolint: lll
func createVirtualNetwork(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
	pp *map[string]interface{},
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	azureConfig, err := getAzureConfig()
	if err != nil {
		return fmt.Errorf("error getting azure config %s", err)
	}
	authorizer, err := getBearerTokenAuthorizer(azureConfig)
	if err != nil {
		return fmt.Errorf("error getting authorizer %s", err)
	}
	virtualNetworksClient := networkSDK.NewVirtualNetworksClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	virtualNetworksClient.Authorizer = authorizer
	virtualNetworkName := uuid.NewV4().String()
	subnetName := "default"
	vnResult, err := virtualNetworksClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		virtualNetworkName,
		networkSDK.VirtualNetwork{
			Location: to.StringPtr("eastus"),
			VirtualNetworkPropertiesFormat: &networkSDK.VirtualNetworkPropertiesFormat{
				AddressSpace: &networkSDK.AddressSpace{
					AddressPrefixes: &[]string{"172.19.0.0/16"},
				},
				Subnets: &[]networkSDK.Subnet{
					{
						Name: &subnetName,
						SubnetPropertiesFormat: &networkSDK.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("172.19.0.0/24"),
						},
					},
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error creating virtual network %s", err)
	}
	if err = vnResult.WaitForCompletion(
		ctx,
		virtualNetworksClient.Client,
	); err != nil {
		return fmt.Errorf("error waiting for virtual network creation complete %s", err)
	}
	subnetID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s",
		azureConfig.SubscriptionID,
		resourceGroup,
		virtualNetworkName,
		subnetName,
	)
	(*pp)["subnetSettings"] = map[string]interface{}{
		"subnetId": subnetID,
	}
	return nil
}

// nolint: lll
func createStorageAccount(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
	pp *map[string]interface{},
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	azureConfig, err := getAzureConfig()
	if err != nil {
		return fmt.Errorf("error getting azure config %s", err)
	}
	authorizer, err := getBearerTokenAuthorizer(azureConfig)
	if err != nil {
		return fmt.Errorf("error getting authorizer %s", err)
	}
	storageClient := storageSDK.NewAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	storageClient.Authorizer = authorizer
	storageAccountName := generate.NewIdentifierOfLength(24)
	result, err := storageClient.Create(
		ctx,
		resourceGroup,
		storageAccountName,
		storageSDK.AccountCreateParameters{
			Sku: &storageSDK.Sku{
				Name: storageSDK.SkuName("Standard_LRS"),
			},
			Kind:     storageSDK.Kind("Storage"),
			Location: to.StringPtr("eastus"),
		},
	)
	if err != nil {
		return fmt.Errorf("error creating storage account %s", err)
	}
	if err := result.WaitForCompletion(
		ctx,
		storageClient.Client,
	); err != nil {
		return fmt.Errorf("error waiting for storage account creation complete %s", err)
	}
	keys, err := storageClient.ListKeys(
		ctx,
		resourceGroup,
		storageAccountName,
	)
	if err != nil {
		return fmt.Errorf("error retrieving storage account credential %s", err)
	}
	primaryKey := *((*keys.Keys)[0].Value)
	connectionString := fmt.Sprintf(
		"DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net",
		storageAccountName,
		primaryKey,
	)
	(*pp)["redisConfiguration"] = make(map[string]interface{})
	configMap := (*pp)["redisConfiguration"].(map[string]interface{})
	configMap["rdb-backup-enabled"] = "enabled"
	configMap["rdb-backup-frequency"] = 60
	configMap["rdb-storage-connection-string"] = connectionString
	return nil
}
