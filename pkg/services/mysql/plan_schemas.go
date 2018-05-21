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
	tierName                string
	tierShortName           string
	allowedHardware         []string
	allowedCores            []int64
	defaultCores            int64
	maxStorage              int64
	allowedBackupRedundancy []string
}

func (t *tierDetails) getSku(pp dbmsProvisioningParameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf(
		"%s_%s_%d",
		t.tierShortName,
		getHardwareFamily(pp),
		t.getCores(pp),
	)
	return sku
}

func generateProvisioningParamsSchema(
	td tierDetails,
) service.InputParametersSchema {
	ips := generateUpdatingParamsSchema(td)
	ips.PropertySchemas["hardwareFamily"] = &service.StringPropertySchema{
		Description:   "Specifies the compute generation to use for the DBMS",
		AllowedValues: td.allowedHardware,
		DefaultValue:  gen5ParamString,
	}
	ips.PropertySchemas["backupRedundancy"] = &service.StringPropertySchema{
		Description:   "Specifies the backup redundancy",
		AllowedValues: td.allowedBackupRedundancy,
		DefaultValue:  "local",
	}
	return *ips
}

func generateUpdatingParamsSchema(
	td tierDetails,
) *service.InputParametersSchema {
	return &service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"cores": &service.IntPropertySchema{
				Description: "Specifies vCores, which represent the logical " +
					"CPU of the underlying hardware",
				AllowedValues: td.allowedCores,
				DefaultValue:  ptr.ToInt64(td.defaultCores),
			},
			"storage": &service.IntPropertySchema{
				Description:  "Specifies the storage in GBs",
				DefaultValue: ptr.ToInt64(10),
				MinValue:     ptr.ToInt64(5),
				MaxValue:     ptr.ToInt64(td.maxStorage),
			},
			"backupRetention": &service.IntPropertySchema{
				Description:  "Specifies the number of days for backup retention",
				DefaultValue: ptr.ToInt64(7),
				MinValue:     ptr.ToInt64(7),
				MaxValue:     ptr.ToInt64(35),
			},
			"sslEnforcement": &service.StringPropertySchema{
				Description: "Specifies whether the server requires the use of TLS" +
					" when connecting. Left unspecified, SSL will be enforced",
				AllowedValues: []string{enabledParamString, disabledParamString},
				DefaultValue:  enabledParamString,
			},
			"firewallRules": &service.ArrayPropertySchema{
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
			},
		},
	}
}

func (t *tierDetails) getCores(pp dbmsProvisioningParameters) int64 {
	if pp.Cores == nil {
		return t.defaultCores
	}
	return *pp.Cores
}

func getStorage(pp dbmsProvisioningParameters) int64 {
	if pp.Storage == nil {
		return 10
	}
	return *pp.Storage
}

func getBackupRetention(pp dbmsProvisioningParameters) int64 {
	if pp.BackupRetention == nil {
		return 7
	}
	return *pp.BackupRetention
}

func isGeoRedundentBackup(pp dbmsProvisioningParameters) bool {
	return pp.BackupRedundancy == "geo"
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
