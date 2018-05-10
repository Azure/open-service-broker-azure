package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func buildGoTemplateParameters(
	svc service.Service,
	provisioningParameters dbmsProvisioningParams,
) map[string]interface{} {
	p := map[string]interface{}{}
	p["version"] = svc.GetProperties().Extended["version"]
	// Only include these if they are not empty.
	// ARM Deployer will fail if the values included are not
	// valid IPV4 addresses (i.e. empty string wil fail)
	if len(provisioningParameters.FirewallRules) > 0 {
		p["firewallRules"] = provisioningParameters.FirewallRules
	} else {
		// Build the azure default
		p["firewallRules"] = []firewallRule{
			{
				Name:    "AllowAzure",
				StartIP: "0.0.0.0",
				EndIP:   "0.0.0.0",
			},
		}
	}
	return p
}
