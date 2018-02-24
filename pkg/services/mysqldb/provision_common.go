package mysqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

func validateServerProvisionParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters " +
				"as *mysql.ServerProvisioningParameters",
		)
	}
	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	if sslEnforcement != "" && sslEnforcement != enabled &&
		sslEnforcement != disabled {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid option: "%s"`, pp.SSLEnforcement),
		)
	}
	for _, firewallRule := range pp.FirewallRules {
		if firewallRule.FirewallRuleName == "" {
			return service.NewValidationError(
				"firewallRuleName",
				"must be set",
			)
		}
		if firewallRule.FirewallIPStart != "" || firewallRule.FirewallIPEnd != "" {
			if firewallRule.FirewallIPStart == "" {
				return service.NewValidationError(
					"firewallStartIPAddress",
					"must be set when firewallEndIPAddress is set",
				)
			}
			if firewallRule.FirewallIPEnd == "" {
				return service.NewValidationError(
					"firewallEndIPAddress",
					"must be set when firewallStartIPAddress is set",
				)
			}
		}
		startIP := net.ParseIP(firewallRule.FirewallIPStart)
		if firewallRule.FirewallIPStart != "" && startIP == nil {
			return service.NewValidationError(
				"firewallStartIPAddress",
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.FirewallIPStart),
			)
		}
		endIP := net.ParseIP(firewallRule.FirewallIPEnd)
		if firewallRule.FirewallIPEnd != "" && endIP == nil {
			return service.NewValidationError(
				"firewallEndIPAddress",
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.FirewallIPEnd),
			)
		}
		//The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
		//Once converted,comparing two IP addresses can be done by using the
		//bytes. Compare function. Per the ARM template documentation,
		//startIP must be <= endIP.
		startBytes := startIP.To4()
		endBytes := endIP.To4()
		if bytes.Compare(startBytes, endBytes) > 0 {
			return service.NewValidationError(
				"firewallEndIPAddress",
				fmt.Sprintf(`invalid value: "%s". must be 
				greater than or equal to firewallStartIPAddress`, firewallRule.FirewallIPEnd),
			)
		}
	}
	return nil
}

func buildGoTemplateParameters(
	provisioningParameters *ServerProvisioningParameters,
) map[string]interface{} {
	p := map[string]interface{}{}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
	if len(provisioningParameters.FirewallRules) > 0 {
		p["firewallRules"] = provisioningParameters.FirewallRules
	} else {
		//Build the azure default
		p["firewallRules"] = []FirewallRule{
			FirewallRule{
				FirewallRuleName: "AllowAzure",
				FirewallIPStart:  "0.0.0.0",
				FirewallIPEnd:    "0.0.0.0",
			},
		}
	}
	return p
}

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			mysqlSDK.NameAvailabilityRequest{
				Name: &serverName,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if *nameAvailability.NameAvailable {
			return serverName, nil
		}
	}
}
