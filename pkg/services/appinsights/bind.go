package appinsights

import (
	"context"

	appInsightsSDK "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"	
)

func (s *serviceManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("Error getting Azure Config: %+v", err)
	}
	subID := azureConfig.SubscriptionID
	resourceGroup := pp.GetString("resourceGroup")
	appInsightsName := dt.AppInsightsName
	apiKeyName := uuid.NewV4().String()
	appInsightsResourceID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s",
		subID,
		resourceGroup,
		appInsightsName,
	)
	apiKeyProperties := appInsightsSDK.APIKeyRequest{
		Name:                  &apiKeyName,
		LinkedReadProperties:  string[]{appInsightsResourceID + "/api", appInsightsResourceID + "/agentconfig"},
		LinkedWriteProperties: string[]{appInsightsResourceID + "/annotations"},
	}		
	result, err := client.Create(context.Background(), resourceGroup, appInsightsName, apiKeyProperties)
	if err != nil {
		return nil, fmt.Errorf("Error creating Application Insights API key %s: %+v", apiKeyName, err)
	}
	if result.APIKey == nil {
		return nil, fmt.Errorf("Error creating Application Insights API key %s: got empty API key", apiKeyName)
	}
	return &bindingDetails{
		APIKeyName: apiKeyName,
		APIKey: result.APIKey,
	}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)
	bd := binding.Details.(*bindingDetails)
	return credentials{
		InstrumentationKey: string(dt.InstrumentationKey),
		AppInsightsName: dt.AppInsightsName,
		APIKey: string(bd.APIKey),
	}, nil
}
