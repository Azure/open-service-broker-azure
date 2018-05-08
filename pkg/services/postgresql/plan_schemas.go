package postgresql

import (
	"fmt"

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
	allowedCores            []int
	defaultCores            int
	tier                    string
	minStorage              int
	maxStorage              int
	defaultStorage          int
	allowedBackupRedundancy []string
	defaultBackupRedundancy string
	minBackupRetention      int
	maxBackupRetention      int
	defaultBackupRetention  int
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

// TODO: krancour: I think it's evident from the code below that we could and
// should generalize this into a framework whereby the schema for provisioning,
// updating, binding params, etc. is codified in some meaninful way and JSON
// schema (included in the catalog for each plan), validation, and derivation
// of default values are ALL driven off of that data.
func generateDBMSPlanSchema(
	schema planSchema,
	includeDBParams bool,
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
			Properties: map[string]service.PropertySchema{
				"name": &service.SimplePropertySchema{
					Type:        "string",
					Description: "Name of firewall rule",
				},
				"startIPAddress": &service.SimplePropertySchema{
					Type:        "string",
					Description: "Start of firewall rule range",
				},
				"endIPAddress": &service.SimplePropertySchema{
					Type:        "string",
					Description: "End of firewall rule range",
				},
			},
		},
	}
	if len(schema.allowedSSLEnforcement) > 1 {
		ps["sslEnforcement"] = &service.SimplePropertySchema{
			Type: "string",
			Description: "Specifies whether the server requires the use of TLS" +
				" when connecting. Left unspecified, SSL will be enforced",
			AllowedValues: schema.allowedSSLEnforcement,
			Default:       schema.defaultSSLEnforcement,
		}
	}
	if len(schema.allowedHardware) > 1 {
		ps["hardwareFamily"] = &service.SimplePropertySchema{
			Type:          "string",
			Description:   "Specifies the compute generation to use for the DBMS",
			AllowedValues: schema.allowedHardware,
			Default:       schema.defaultHardware,
		}
	}
	if len(schema.allowedCores) > 1 {
		ps["cores"] = &service.SimplePropertySchema{
			Type: "number",
			Description: "Specifies vCores, which represent the logical " +
				"CPU of the underlying hardware",
			AllowedValues: schema.allowedCores,
			Default:       schema.defaultCores,
		}
	}
	if schema.maxStorage > schema.minStorage {
		ps["storage"] = &service.NumericPropertySchema{
			Type:        "number",
			Description: "Specifies the storage in GBs",
			Default:     schema.defaultStorage,
			Minimum:     schema.minStorage,
			Maximum:     schema.maxStorage,
		}
	}
	if schema.maxBackupRetention > schema.minBackupRetention {
		ps["backupRetention"] = &service.NumericPropertySchema{
			Type:        "number",
			Description: "Specifies the number of days for backup retention",
			Default:     schema.minBackupRetention,
			Minimum:     schema.minBackupRetention,
			Maximum:     schema.maxBackupRetention,
		}
	}
	if len(schema.allowedBackupRedundancy) > 1 {
		ps["backupRedundancy"] = &service.SimplePropertySchema{
			Type:          "string",
			Description:   "Specifies the backup redundancy",
			AllowedValues: schema.allowedBackupRedundancy,
			Default:       schema.defaultBackupRedundancy,
		}
	}
	if includeDBParams {
		ps["extensions"] = dbExtensionsSchema
	}
	return service.InputParametersSchema{
		Properties: ps,
	}
}

func (p *planSchema) getCores(pp dbmsProvisioningParameters) int {
	// If you get a choice and you've made a choice...
	if len(p.allowedCores) > 1 && pp.Cores != nil {
		return *pp.Cores
	}
	return p.defaultCores
}

func (p *planSchema) getStorage(pp dbmsProvisioningParameters) int {
	// If you get a choice and you've made a choice...
	if p.maxStorage > p.minStorage && pp.Storage != nil {
		return *pp.Storage
	}
	return p.defaultStorage
}

func (p *planSchema) getBackupRetention(pp dbmsProvisioningParameters) int {
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
