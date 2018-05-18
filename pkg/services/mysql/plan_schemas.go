package mysql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const (
	gen4TemplateString = "Gen4"
	gen5TemplateString = "Gen5"
	gen4ParamString    = "gen4"
	gen5ParamString    = "gen5"
)

type tierDetails struct {
	allowedHardware         []string
	allowedCores            []int64
	defaultCores            int64
	tier                    string
	minStorage              int64
	maxStorage              int64
	defaultStorage          int64
	allowedBackupRedundancy []string
	defaultBackupRedundancy string
	minBackupRetention      int64
	maxBackupRetention      int64
	defaultBackupRetention  int64
}

func (t *tierDetails) getSku(pp dbmsProvisioningParameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf(
		"%s_%s_%d",
		t.tier,
		getHardwareFamily(pp),
		t.getCores(pp),
	)
	return sku
}

func generateDBMSPlanSchema(
	td tierDetails,
) service.InputParametersSchema {
	ps := map[string]service.PropertySchema{}
	ps["firewallRules"] = &service.ArrayPropertySchema{
		Description: "Firewall rules to apply to instance. " +
			"If left unspecified, defaults to only Azure IPs",
		ItemsSchema: &service.ObjectPropertySchema{
			Description: "Individual Firewall Rule",
			RequiredProperties: []string{
				"name",
				"startIPAddress",
				"endIPAddress",
			},
			PropertySchemas: map[string]service.PropertySchema{
				"name": &service.StringPropertySchema{
					Description: "Name of firewall rule",
				},
				"startIPAddress": &service.StringPropertySchema{
					Description:             "Start of firewall rule range",
					CustomPropertyValidator: ipValidator,
				},
				"endIPAddress": &service.StringPropertySchema{
					Description:             "End of firewall rule range",
					CustomPropertyValidator: ipValidator,
				},
			},
			CustomPropertyValidator: firewallRuleValidator,
		},
		DefaultValue: []interface{}{
			map[string]interface{}{
				"name":           "AllowAzure",
				"startIPAddress": "0.0.0.0",
				"endIPAddress":   "0.0.0.0",
			},
		},
	}
	ps["sslEnforcement"] = &service.StringPropertySchema{
		Description: "Specifies whether the server requires the use of TLS" +
			" when connecting. Left unspecified, SSL will be enforced",
		AllowedValues: []string{enabledParamString, disabledParamString},
		DefaultValue:  enabledParamString,
	}
	if len(td.allowedHardware) > 1 {
		ps["hardwareFamily"] = &service.StringPropertySchema{
			Description:   "Specifies the compute generation to use for the DBMS",
			AllowedValues: td.allowedHardware,
			DefaultValue:  gen5ParamString,
		}
	}
	if len(td.allowedCores) > 1 {
		ps["cores"] = &service.IntPropertySchema{
			Description: "Specifies vCores, which represent the logical " +
				"CPU of the underlying hardware",
			AllowedValues: td.allowedCores,
			DefaultValue:  ptr.ToInt64(td.defaultCores),
		}
	}
	if td.maxStorage > td.minStorage {
		ps["storage"] = &service.IntPropertySchema{
			Description:  "Specifies the storage in GBs",
			DefaultValue: ptr.ToInt64(td.defaultStorage),
			MinValue:     ptr.ToInt64(td.minStorage),
			MaxValue:     ptr.ToInt64(td.maxStorage),
		}
	}
	if td.maxBackupRetention > td.minBackupRetention {
		ps["backupRetention"] = &service.IntPropertySchema{
			Description:  "Specifies the number of days for backup retention",
			DefaultValue: ptr.ToInt64(td.minBackupRetention),
			MinValue:     ptr.ToInt64(td.minBackupRetention),
			MaxValue:     ptr.ToInt64(td.maxBackupRetention),
		}
	}
	if len(td.allowedBackupRedundancy) > 1 {
		ps["backupRedundancy"] = &service.StringPropertySchema{
			Description:   "Specifies the backup redundancy",
			AllowedValues: td.allowedBackupRedundancy,
			DefaultValue:  td.defaultBackupRedundancy,
		}
	}
	return service.InputParametersSchema{
		PropertySchemas: ps,
	}
}

func (t *tierDetails) getCores(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if len(t.allowedCores) > 1 && pp.Cores != nil {
		return *pp.Cores
	}
	return t.defaultCores
}

func (t *tierDetails) getStorage(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if t.maxStorage > t.minStorage && pp.Storage != nil {
		return *pp.Storage
	}
	return t.defaultStorage
}

func (t *tierDetails) getBackupRetention(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if t.maxBackupRetention > t.minBackupRetention && pp.BackupRetention != nil {
		return *pp.BackupRetention
	}
	return t.defaultBackupRetention
}

func (t *tierDetails) isGeoRedundentBackup(pp dbmsProvisioningParameters) bool {
	// If you get a choice and you've made a choice...
	if len(t.allowedBackupRedundancy) > 1 && pp.BackupRedundancy != "" {
		return pp.BackupRedundancy == "geo"
	}
	return t.defaultBackupRedundancy == "geo"
}

func getHardwareFamily(pp dbmsProvisioningParameters) string {
	if pp.HardwareFamily == "" {
		return gen5TemplateString
	}
	if pp.HardwareFamily == gen4ParamString {
		return gen4TemplateString
	}
	return gen5TemplateString
}

func isSSLRequired(pp dbmsProvisioningParameters) bool {
	return pp.SSLEnforcement != disabledParamString
}

func getFirewallRules(
	pp dbmsProvisioningParameters,
) []firewallRule {
	if len(pp.FirewallRules) > 0 {
		return pp.FirewallRules
	}
	return []firewallRule{
		{
			Name:    "AllowAzure",
			StartIP: "0.0.0.0",
			EndIP:   "0.0.0.0",
		},
	}
}

func ipValidator(context, value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return service.NewValidationError(
			context,
			fmt.Sprintf(`"%s" is not a valid IP address`, value),
		)
	}
	return nil
}

func firewallRuleValidator(
	context string,
	valMap map[string]interface{},
) error {
	startIP := net.ParseIP(valMap["startIPAddress"].(string))
	endIP := net.ParseIP(valMap["endIPAddress"].(string))
	// The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
	// Once converted,comparing two IP addresses can be done by using the
	// bytes. Compare function. Per the ARM template documentation,
	// startIP must be <= endIP.
	startBytes := startIP.To4()
	endBytes := endIP.To4()
	if bytes.Compare(startBytes, endBytes) > 0 {
		return service.NewValidationError(
			context,
			fmt.Sprintf(
				`endIPAddress "%s" is not greater than or equal to startIPAddress "%s"`,
				endIP,
				startIP,
			),
		)
	}
	return nil
}
