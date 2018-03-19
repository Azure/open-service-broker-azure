package mssql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func validateDBMSProvisionParameters(pp dbmsProvisioningParams) error {
	for _, firewallRule := range pp.FirewallRules {
		if firewallRule.Name == "" {
			return service.NewValidationError(
				"name",
				"must be set",
			)
		}
		if firewallRule.StartIP != "" ||
			firewallRule.EndIP != "" {
			if firewallRule.StartIP == "" {
				return service.NewValidationError(
					"startIPAddress",
					"must be set when endIPAddress is set",
				)
			}
			if firewallRule.EndIP == "" {
				return service.NewValidationError(
					"endIPAddress",
					"must be set when startIPAddress is set",
				)
			}
		}
		startIP := net.ParseIP(firewallRule.StartIP)
		if firewallRule.StartIP != "" && startIP == nil {
			return service.NewValidationError(
				"startIPAddress",
				fmt.Sprintf(
					`invalid value: "%s"`,
					firewallRule.StartIP,
				),
			)
		}
		endIP := net.ParseIP(firewallRule.EndIP)
		if firewallRule.EndIP != "" && endIP == nil {
			return service.NewValidationError(
				"endIPAddress",
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.EndIP),
			)
		}
		// The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
		// Once converted,comparing two IP addresses can be done by using the
		// bytes. Compare function. Per the ARM template documentation,
		// startIP must be <= endIP.
		startBytes := startIP.To4()
		endBytes := endIP.To4()
		if bytes.Compare(startBytes, endBytes) > 0 {
			return service.NewValidationError(
				"endIPAddress",
				fmt.Sprintf(`invalid value: "%s". must be 
				greater than or equal to startIPAddress`,
					firewallRule.EndIP,
				),
			)
		}
	}
	return nil
}

func buildGoTemplateParameters(
	provisioningParameters dbmsProvisioningParams,
) map[string]interface{} {
	p := map[string]interface{}{}
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
