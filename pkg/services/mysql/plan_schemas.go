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

type planSchema struct {
	defaultFirewallRules    []firewallRule
	allowedSSLEnforcement   []string
	defaultSSLEnforcement   string
	defaultHardware         string
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

func (p *planSchema) getSku(pp dbmsProvisioningParameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf(
		"%s_%s_%d",
		p.tier,
		p.getHardwareFamily(pp),
		p.getCores(pp),
	)
	return sku
}

func generateDBMSPlanSchema(
	schema planSchema,
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
	if len(schema.allowedSSLEnforcement) > 1 {
		ps["sslEnforcement"] = &service.StringPropertySchema{
			Description: "Specifies whether the server requires the use of TLS" +
				" when connecting. Left unspecified, SSL will be enforced",
			AllowedValues: schema.allowedSSLEnforcement,
			DefaultValue:  schema.defaultSSLEnforcement,
		}
	}
	if len(schema.allowedHardware) > 1 {
		ps["hardwareFamily"] = &service.StringPropertySchema{
			Description:   "Specifies the compute generation to use for the DBMS",
			AllowedValues: schema.allowedHardware,
			DefaultValue:  schema.defaultHardware,
		}
	}
	if len(schema.allowedCores) > 1 {
		ps["cores"] = &service.IntPropertySchema{
			Description: "Specifies vCores, which represent the logical " +
				"CPU of the underlying hardware",
			AllowedValues: schema.allowedCores,
			DefaultValue:  ptr.ToInt64(schema.defaultCores),
		}
	}
	if schema.maxStorage > schema.minStorage {
		ps["storage"] = &service.IntPropertySchema{
			Description:  "Specifies the storage in GBs",
			DefaultValue: ptr.ToInt64(schema.defaultStorage),
			MinValue:     ptr.ToInt64(schema.minStorage),
			MaxValue:     ptr.ToInt64(schema.maxStorage),
		}
	}
	if schema.maxBackupRetention > schema.minBackupRetention {
		ps["backupRetention"] = &service.IntPropertySchema{
			Description:  "Specifies the number of days for backup retention",
			DefaultValue: ptr.ToInt64(schema.minBackupRetention),
			MinValue:     ptr.ToInt64(schema.minBackupRetention),
			MaxValue:     ptr.ToInt64(schema.maxBackupRetention),
		}
	}
	if len(schema.allowedBackupRedundancy) > 1 {
		ps["backupRedundancy"] = &service.StringPropertySchema{
			Description:   "Specifies the backup redundancy",
			AllowedValues: schema.allowedBackupRedundancy,
			DefaultValue:  schema.defaultBackupRedundancy,
		}
	}
	return service.InputParametersSchema{
		PropertySchemas: ps,
	}
}

func (p *planSchema) getCores(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if len(p.allowedCores) > 1 && pp.Cores != nil {
		return *pp.Cores
	}
	return p.defaultCores
}

func (p *planSchema) getStorage(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if p.maxStorage > p.minStorage && pp.Storage != nil {
		return *pp.Storage
	}
	return p.defaultStorage
}

func (p *planSchema) getBackupRetention(pp dbmsProvisioningParameters) int64 {
	// If you get a choice and you've made a choice...
	if p.maxBackupRetention > p.minBackupRetention && pp.BackupRetention != nil {
		return *pp.BackupRetention
	}
	return p.defaultBackupRetention
}

func (p *planSchema) isGeoRedundentBackup(pp dbmsProvisioningParameters) bool {
	// If you get a choice and you've made a choice...
	if len(p.allowedBackupRedundancy) > 1 && pp.BackupRedundancy != "" {
		return pp.BackupRedundancy == "geo"
	}
	return p.defaultBackupRedundancy == "geo"
}

func (p *planSchema) getHardwareFamily(pp dbmsProvisioningParameters) string {
	var hardwareSelection string
	// If you get a choice and you've made a choice...
	if len(p.allowedHardware) > 1 && hardwareSelection == "" {
		hardwareSelection = pp.HardwareFamily
	} else {
		hardwareSelection = p.defaultHardware
	}
	// Translate to a value usable in the ARM templates.
	// TODO: It might be better for this object not to know so much about how it
	// is ultimately used-- i.e. ARM-template-awareness.
	if hardwareSelection == gen4ParamString {
		return gen4TemplateString
	}
	return gen5TemplateString
}

func (p *planSchema) isSSLRequired(pp dbmsProvisioningParameters) bool {
	// If you get a choice and you've made a choice...
	if len(p.allowedSSLEnforcement) > 1 && pp.SSLEnforcement != "" {
		return pp.SSLEnforcement == enabledParamString
	}
	return p.defaultSSLEnforcement == enabledParamString
}

func (p *planSchema) getFirewallRules(
	pp dbmsProvisioningParameters,
) []firewallRule {
	if len(pp.FirewallRules) > 0 {
		return pp.FirewallRules
	}
	return p.defaultFirewallRules
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
