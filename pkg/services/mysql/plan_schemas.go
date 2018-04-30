package mysql

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
	validCores              []int
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

func (p *planSchema) buildSku(
	pp dbmsProvisioningParameters,
) (string, error) {
	hardwareFamily, err := p.generateHardwareFamilyString(pp)
	if err != nil {
		return "", fmt.Errorf("error building sku: %s", err)
	}
	var cores int

	if pp.Cores == nil {
		cores = p.defaultCores
	} else {
		cores = *pp.Cores
	}
	//The name of the sku, typically:
	//tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf("%s_%s_%d", p.tier, hardwareFamily, cores)
	return sku, nil
}

func (p *planSchema) generateHardwareFamilyString(
	pp dbmsProvisioningParameters,
) (string, error) {
	if pp.HardwareFamily == "" {
		return gen5TemplateString, nil
	} else if pp.HardwareFamily == gen4ParamString {
		return gen4TemplateString, nil
	} else if pp.HardwareFamily == gen5ParamString {
		return gen5TemplateString, nil
	}
	return "", fmt.Errorf("unknown hardware family")
}

func generateDBMSPlanSchema(
	schema planSchema,
) map[string]service.ParameterSchema {
	return map[string]service.ParameterSchema{
		"firewallRules": &service.ArrayParameterSchema{
			Description: "Firewall rules to apply to instance. " +
				"If left unspecified, defaults to only Azure IPs",
			ItemsSchema: &service.ObjectParameterSchema{
				Description: "Individual Firewall Rule",
				Properties: map[string]service.ParameterSchema{
					"name": &service.SimpleParameterSchema{
						Type:        "string",
						Description: "Name of firewall rule",
						Required:    true,
					},
					"startIPAddress": &service.SimpleParameterSchema{
						Type:        "string",
						Description: "Start of firewall rule range",
						Required:    true,
					},
					"endIPAddress": &service.SimpleParameterSchema{
						Type:        "string",
						Description: "End of firewall rule range",
						Required:    true,
					},
				},
			},
		},
		"sslEnforcement": &service.SimpleParameterSchema{
			Type: "string",
			Description: "Specifies whether the server requires the use of TLS" +
				" when connecting. Left unspecified, SSL will be enforced",
			AllowedValues: []string{"enabled", "disabled"},
			Default:       "enabled",
		},
		"hardwareFamily": &service.SimpleParameterSchema{
			Type:          "string",
			Description:   "Specifies the compute generation to use for the DBMS",
			AllowedValues: schema.allowedHardware,
			Default:       schema.defaultHardware,
		},
		"cores": &service.SimpleParameterSchema{
			Type: "number",
			Description: "Specifies vCores, which represent the logical CPU " +
				"of the underlying hardware",
			AllowedValues: schema.validCores,
			Default:       schema.defaultCores,
		},
		"storage": &service.NumericParameterSchema{
			Type:        "number",
			Description: "Specifies the storage in GBs",
			Default:     schema.defaultStorage,
			Minimum:     schema.minStorage,
			Maximum:     schema.maxStorage,
		},
		"backupRetention": &service.NumericParameterSchema{
			Type:        "number",
			Description: "Specifies the number of days for backup retention",
			Default:     schema.minBackupRetention,
			Minimum:     schema.minBackupRetention,
			Maximum:     schema.maxBackupRetention,
		},
		"backupRedundancy": &service.SimpleParameterSchema{
			Type:          "string",
			Description:   "Specifies the backup redundancy",
			AllowedValues: schema.allowedBackupRedundancy,
			Default:       schema.defaultBackupRedundancy,
		},
	}
}

func (p *planSchema) getCores(pp dbmsProvisioningParameters) int {
	if pp.Cores != nil {
		return *pp.Cores
	}
	return p.defaultCores

}

func (p *planSchema) getStorage(pp dbmsProvisioningParameters) int {
	if pp.Storage != nil {
		return *pp.Storage
	}
	return p.defaultStorage

}

func (p *planSchema) getBackupRetention(pp dbmsProvisioningParameters) int {
	if pp.BackupRetention != nil {
		return *pp.BackupRetention
	}
	return p.defaultBackupRetention

}

func (p *planSchema) isGeoRedundentBackup(pp dbmsProvisioningParameters) bool {
	return pp.BackupRedundancy == "geo"
}

func (p *planSchema) getHardwareFamily(pp dbmsProvisioningParameters) string {
	if pp.HardwareFamily == "" {
		if p.defaultHardware == gen4ParamString {
			return gen4TemplateString
		}
		return gen5TemplateString
	} else if pp.HardwareFamily == gen4ParamString {
		return gen4TemplateString
	}
	return gen5TemplateString
}

func (p *planSchema) isSSLRequired(pp dbmsProvisioningParameters) bool {
	if pp.SSLEnforcement != "" {
		return pp.SSLEnforcement == enabledParamString
	}
	return true
}

func (p *planSchema) getFirewallRules(
	pp dbmsProvisioningParameters,
) []firewallRule {
	if len(pp.FirewallRules) > 0 {
		return pp.FirewallRules
	}
	return p.defaultFirewallRules
}
