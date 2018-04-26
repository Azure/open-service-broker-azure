package mysql

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabled           = "enabled"
	disabled          = "disabled"
	enabledARMString  = "Enabled"
	disabledARMString = "Disabled"
)

func validateDBMSProvisionParameters(
	plan service.Plan,
	pp dbmsProvisioningParameters,
) error {
	if plan == nil {
		return fmt.Errorf("plan invalid")
	}
	s, ok := plan.GetProperties().Extended["provisionSchema"]
	if !ok {
		return fmt.Errorf("invalid plan, schema not found")
	}
	schema := s.(planSchema)
	if err := schema.validateProvisionParameters(pp); err != nil {
		return err
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
		if firewallRule.Name == "" {
			return service.NewValidationError(
				"name",
				"must be set",
			)
		}
		if firewallRule.StartIP != "" || firewallRule.EndIP != "" {
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
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.StartIP),
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
				greater than or equal tostartIPAddress`,
					firewallRule.EndIP),
			)
		}
	}
	return nil
}

func buildGoTemplateParameters(
	instance service.Instance,
) (map[string]interface{}, error) {

	plan := instance.Plan

	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	pp := dbmsProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, err
	}

	p := map[string]interface{}{}
	schema := plan.GetProperties().Extended["provisionSchema"].(planSchema)

	sku, err := schema.buildSku(pp)
	if err != nil {
		return nil, err
	}
	p["sku"] = sku
	p["tier"] = plan.GetProperties().Extended["tier"]
	p["cores"] = schema.getCores(pp)
	p["storage"] = schema.getStorage(pp) * 1024 //storage is in MB to arm :/
	p["backupRetention"] = schema.getBackupRetention(pp)
	p["hardwareFamily"] = schema.getHardwareFamily(pp)
	if schema.isGeoRedundentBackup(pp) {
		p["geoRedundantBackup"] = enabledARMString
	}

	p["serverName"] = dt.ServerName
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	if dt.EnforceSSL {
		p["sslEnforcement"] = enabledARMString
	} else {
		p["sslEnforcement"] = disabledARMString
	}

	p["version"] = instance.Service.GetProperties().Extended["version"]

	// Only include these if they are not empty.
	// ARM Deployer will fail if the values included are not
	// valid IPV4 addresses (i.e. empty string wil fail)
	if len(pp.FirewallRules) > 0 {
		p["firewallRules"] = pp.FirewallRules
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
	return p, nil
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
