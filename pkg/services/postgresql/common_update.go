package postgresql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/slice"
	"github.com/Azure/open-service-broker-azure/pkg/types"
	_ "github.com/lib/pq" // Postgres SQL driver
)

func validateDBMSUpdateParameters(
	plan service.Plan,
	oldPP dbmsProvisioningParameters,
	up dbmsUpdatingParameters,
) error {
	if plan == nil {
		return fmt.Errorf("plan invalid")
	}
	s, ok := plan.GetProperties().Extended["updateSchema"]
	if !ok {
		return fmt.Errorf("invalid plan, schema not found")
	}
	schema := s.(planSchema)
	if up.SSLEnforcement != "" &&
		!slice.ContainsString(schema.allowedSSLEnforcement, up.SSLEnforcement) {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid option: "%s"`, up.SSLEnforcement),
		)
	}
	for _, firewallRule := range up.FirewallRules {
		if firewallRule.Name == "" {
			return service.NewValidationError(
				"ruleName",
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
		endIP := net.ParseIP(firewallRule.StartIP)
		if firewallRule.EndIP != "" && endIP == nil {
			return service.NewValidationError(
				"endIPAddress",
				fmt.Sprintf(
					`invalid value: "%s"`,
					firewallRule.EndIP,
				),
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
					firewallRule.EndIP),
			)
		}
	}

	// cores
	if up.Cores != nil && !slice.ContainsInt(schema.allowedCores, *up.Cores) {
		return service.NewValidationError(
			"cores",
			fmt.Sprintf(`invalid value: "%d"`, *up.Cores),
		)
	}

	// storage
	if up.Storage != nil {
		if *up.Storage < schema.minStorage || *up.Storage > schema.maxStorage {
			return service.NewValidationError(
				"storage",
				fmt.Sprintf(`invalid value: "%d"`, *up.Storage),
			)
		}
		// This is the only real domain specific thing here that needs to be
		// validated against the provisioned instance
		if oldPP.Storage != nil && *up.Storage < *oldPP.Storage {
			return service.NewValidationError(
				"storage",
				fmt.Sprintf(`invalid value: "%d". cannot reduce storage`, *up.Storage),
			)
		}

	}

	// backupRetation
	if up.BackupRetention != nil &&
		(*up.BackupRetention < schema.minBackupRetention ||
			*up.BackupRetention > schema.maxBackupRetention) {
		return service.NewValidationError(
			"backupRetention",
			fmt.Sprintf(`invalid value: "%d"`, *up.BackupRetention),
		)
	}

	return nil
}

func buildGoUpdateTemplateParameters(
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
	mergedParams := instance.ProvisioningParameters
	// create a copy of the provision parameters merged with updating parameters
	// and use this to generate the values that follow via the schema. Only
	// overwrite the provision params that have values in the update params
	for key, value := range instance.UpdatingParameters {
		if !types.IsEmpty(value) {
			mergedParams[key] = value
		}

	}
	if err := service.GetStructFromMap(mergedParams, &pp); err != nil {
		return nil, err
	}

	schema := plan.GetProperties().Extended["provisionSchema"].(planSchema)

	p := map[string]interface{}{}
	p["sku"] = schema.getSku(pp)
	p["tier"] = plan.GetProperties().Extended["tier"]
	p["cores"] = schema.getCores(pp)
	p["storage"] = schema.getStorage(pp) * 1024 //storage is in MB to arm :/
	p["backupRetention"] = schema.getBackupRetention(pp)
	p["hardwareFamily"] = schema.getHardwareFamily(pp)
	if schema.isGeoRedundentBackup(pp) {
		p["geoRedundantBackup"] = enabledARMString
	}
	p["version"] = instance.Service.GetProperties().Extended["version"]
	p["serverName"] = dt.ServerName
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	if schema.isSSLRequired(pp) {
		p["sslEnforcement"] = enabledARMString
	} else {
		p["sslEnforcement"] = disabledARMString
	}
	p["firewallRules"] = schema.getFirewallRules(pp)

	return p, nil
}
