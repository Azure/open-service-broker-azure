package appinsights

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters
	bd := binding.Details.(*bindingDetails)

	resourceGroup := pp.GetString("resourceGroup")
	appInsightsName := dt.AppInsightsName
	apiKeyName := bd.APIKeyName

	if _, err := client.Delete(
		context.Background(),
		resourceGroup,
		appInsightsName,
		apiKeyName,
	); err != nil {
		return fmt.Errorf("Error deleting Application Insights API key %s: %+v", apiKeyName, err)
	}

	return nil
}
