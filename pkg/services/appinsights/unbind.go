package appinsights

import (
	"context"
	"fmt"

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
	apiKeyID := bd.APIKeyID

	if _, err := s.appInsightsAPIKeyClient.Delete(
		context.Background(),
		resourceGroup,
		appInsightsName,
		apiKeyID,
	); err != nil {
		return fmt.Errorf(
			"Error deleting Application Insights API key %s: %+v",
			apiKeyID,
			err,
		)
	}

	return nil
}
