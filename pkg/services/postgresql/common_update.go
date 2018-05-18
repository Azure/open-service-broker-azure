package postgresql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/slice"
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
	schema := s.(tierDetails)
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
	if up.Cores != nil && !slice.ContainsInt64(schema.allowedCores, *up.Cores) {
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

	// This is the only real domain specific thing here that needs to be
	// validated against the provisioned instance. It is encapsulated in a function
	// so it can be easily refactored once we move to the common validation
	err := validateStorageUpdate(oldPP, up)
	return err
}

func validateStorageUpdate(
	pp dbmsProvisioningParameters,
	up dbmsUpdatingParameters,
) error {
	if up.Storage != nil {
		if pp.Storage != nil && *up.Storage < *pp.Storage {
			return service.NewValidationError(
				"storage",
				fmt.Sprintf(`invalid value: "%d". cannot reduce storage`, *up.Storage),
			)
		}
	}
	return nil
}
