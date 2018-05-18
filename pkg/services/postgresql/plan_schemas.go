package postgresql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planSchema struct {
	defaultFirewallRules    []interface{}
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

func generateDBMSPlanSchema(
	schema planSchema,
	includeDBParams bool,
) *service.InputParametersSchema {
	ps := map[string]service.PropertySchema{
		"version": &service.StringPropertySchema{
			AllowedValues: []string{"9.6"},
			DefaultValue:  "9.6",
		},
		"tier": &service.StringPropertySchema{
			Description:   "Specifies the service tier",
			AllowedValues: []string{schema.tier},
			DefaultValue:  schema.tier,
		},
		"hardwareFamily": &service.StringPropertySchema{
			Description:   "Specifies the compute generation to use for the DBMS",
			AllowedValues: schema.allowedHardware,
			DefaultValue:  schema.defaultHardware,
		},
		"cores": &service.IntPropertySchema{
			Description: "Specifies vCores, which represent the logical " +
				"CPU of the underlying hardware",
			AllowedValues: schema.allowedCores,
			DefaultValue:  ptr.ToInt64(schema.defaultCores),
		},
		"storage": &service.IntPropertySchema{
			Description:  "Specifies the storage in GBs",
			DefaultValue: ptr.ToInt64(schema.defaultStorage),
			MinValue:     ptr.ToInt64(schema.minStorage),
			MaxValue:     ptr.ToInt64(schema.maxStorage),
		},
		"backupRetention": &service.IntPropertySchema{
			Description:  "Specifies the number of days for backup retention",
			DefaultValue: ptr.ToInt64(schema.minBackupRetention),
			MinValue:     ptr.ToInt64(schema.minBackupRetention),
			MaxValue:     ptr.ToInt64(schema.maxBackupRetention),
		},
		"backupRedundancy": &service.StringPropertySchema{
			Description:   "Specifies the backup redundancy",
			AllowedValues: schema.allowedBackupRedundancy,
			DefaultValue:  schema.defaultBackupRedundancy,
		},
		"sslEnforcement": &service.StringPropertySchema{
			Description: "Specifies whether the server requires the use of TLS" +
				" when connecting. Left unspecified, SSL will be enforced",
			AllowedValues: schema.allowedSSLEnforcement,
			DefaultValue:  schema.defaultSSLEnforcement,
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
	}
	if includeDBParams {
		ps["extensions"] = dbExtensionsSchema
	}
	return &service.InputParametersSchema{
		PropertySchemas: ps,
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
