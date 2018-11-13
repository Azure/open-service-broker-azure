package storage

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildGoTemplate(
	instance service.Instance,
	parameter service.ProvisioningParameters,
) map[string]interface{} {
	dt := instance.Details.(*instanceDetails)

	location := parameter.GetString("location")
	nonHTTPSEnabled := parameter.GetString("enableNonHttpsTraffic")
	goTemplateParams := map[string]interface{}{
		"name":                    dt.StorageAccountName,
		"location":                location,
		"supportHttpsTrafficOnly": nonHTTPSEnabled == disabled,
		"accountType":             parameter.GetString("accountType"),
	}

	serviceName := instance.Service.GetName()
	switch serviceName {
	case serviceGeneralPurposeV1:
		goTemplateParams["kind"] = "Storage"
	case serviceGeneralPurposeV2:
		goTemplateParams["kind"] = "StorageV2"
		goTemplateParams["accessTier"] = parameter.GetString("accessTier")
	case serviceBlobAccount, serviceBlobAllInOne:
		goTemplateParams["kind"] = "BlobStorage"
		goTemplateParams["accessTier"] = parameter.GetString("accessTier")
	}
	return goTemplateParams
}
