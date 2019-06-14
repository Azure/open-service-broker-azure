package appinsights

import (
	"context"
	"fmt"
	"strings"

	appInsightsSDK "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
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
	apiKeyName := generate.NewIdentifier()
	appInsightsResourceID := fmt.Sprintf(
		"/subscriptions/%s/"+
			"resourceGroups/%s/"+
			"providers/Microsoft.Insights/components/%s",
		subID,
		resourceGroup,
		appInsightsName,
	)
	apiKeyProperties := appInsightsSDK.APIKeyRequest{
		Name: &apiKeyName,
		LinkedReadProperties: &[]string{
			appInsightsResourceID + "/api",
			appInsightsResourceID + "/agentconfig",
		},
		LinkedWriteProperties: &[]string{
			appInsightsResourceID + "/annotations",
		},
	}
	result, err := s.appInsightsAPIKeyClient.Create(
		context.Background(),
		resourceGroup,
		appInsightsName,
		apiKeyProperties,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Error creating Application Insights API key %s: %+v",
			apiKeyName,
			err,
		)
	}
	if result.APIKey == nil {
		return nil, fmt.Errorf(
			"Error creating Application Insights API key %s: got empty API key",
			apiKeyName,
		)
	}
	if result.ID == nil {
		return nil, fmt.Errorf(
			"Error creating Application Insights API key %s: got empty API key ID",
			apiKeyName,
		)
	}
	// Note that, AppID is different from ApplicationID by design ...
	properties := result.ApplicationInsightsComponentProperties
	if properties == nil {
		return nil, fmt.Errorf(
			"Error creating Application Insights API key %s: got empty properties ID",
			apiKeyName,
		)
		if (*properties).AppID == nil {
			return nil, fmt.Errorf(
				"Error creating Application Insights API key %s: got empty App ID",
				apiKeyName,
			)
		}
	}

	apiKeyID := (*result.ID)[strings.LastIndex(*result.ID, "/")+1:]

	return &bindingDetails{
		AppID:    *(*properties).AppID,
		APIKeyID: apiKeyID,
		APIKey:   service.SecureString(*result.APIKey),
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
		AppID:              bd.AppID,
		APIKey:             string(bd.APIKey),
	}, nil
}
