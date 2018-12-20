package postgresql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/schemas"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type tierDetails struct {
	tierName                string
	tierShortName           string
	allowedCores            []int64
	defaultCores            int64
	maxStorage              int64
	allowedBackupRedundancy []string
}

func (t *tierDetails) getSku(pp service.ProvisioningParameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	sku := fmt.Sprintf(
		"%s_Gen5_%d",
		t.tierShortName,
		pp.GetInt64("cores"),
	)
	return sku
}

func generateProvisioningParamsSchema(
	td tierDetails,
	includeDBParams bool,
) service.InputParametersSchema {
	ips := generateUpdatingParamsSchema(td)
	ips.RequiredProperties = append(ips.RequiredProperties, "location")
	ips.PropertySchemas["location"] = schemas.GetLocationSchema()
	ips.RequiredProperties = append(ips.RequiredProperties, "resourceGroup")
	ips.PropertySchemas["resourceGroup"] = schemas.GetResourceGroupSchema()
	ips.PropertySchemas["backupRedundancy"] = &service.StringPropertySchema{
		Title:         "Backup redundancy",
		Description:   "Specifies the backup redundancy",
		AllowedValues: td.allowedBackupRedundancy,
		DefaultValue:  "local",
	}
	ips.PropertySchemas["tags"] = &service.ObjectPropertySchema{
		Title: "Tags",
		Description: "Tags to be applied to new resources," +
			" specified as key/value pairs.",
		Additional: &service.StringPropertySchema{},
	}
	if includeDBParams {
		ips.PropertySchemas["extensions"] = dbExtensionsSchema
	}
	return ips
}

func generateUpdatingParamsSchema(
	td tierDetails,
) service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"cores": &service.IntPropertySchema{
				Title: "Cores",
				Description: "Specifies vCores, which represent the logical " +
					"CPU of the underlying hardware",
				AllowedValues: td.allowedCores,
				DefaultValue:  ptr.ToInt64(td.defaultCores),
			},
			"storage": &service.IntPropertySchema{
				Title:        "Storage",
				Description:  "Specifies the storage in GBs",
				DefaultValue: ptr.ToInt64(10),
				MinValue:     ptr.ToInt64(5),
				MaxValue:     ptr.ToInt64(td.maxStorage),
			},
			"backupRetention": &service.IntPropertySchema{
				Title:        "Backup retention",
				Description:  "Specifies the number of days for backup retention",
				DefaultValue: ptr.ToInt64(7),
				MinValue:     ptr.ToInt64(7),
				MaxValue:     ptr.ToInt64(35),
			},
			"sslEnforcement": &service.StringPropertySchema{
				Title: "SSL enforcement",
				Description: "Specifies whether the server requires the use of TLS" +
					" when connecting. Left unspecified, SSL will be enforced",
				AllowedValues: []string{enabledParamString, disabledParamString},
				DefaultValue:  enabledParamString,
			},
			"firewallRules": &service.ArrayPropertySchema{
				Title: "Firewall rules",
				Description: "Firewall rules to apply to instance. " +
					"If left unspecified, defaults to only Azure IPs",
				ItemsSchema: &service.ObjectPropertySchema{
					Title:       "Firewall rule",
					Description: "Individual Firewall Rule",
					RequiredProperties: []string{
						"name",
						"startIPAddress",
						"endIPAddress",
					},
					PropertySchemas: map[string]service.PropertySchema{
						"name": &service.StringPropertySchema{
							Title:       "Name",
							Description: "Name of firewall rule",
						},
						"startIPAddress": &service.StringPropertySchema{
							Title:                   "Start IP address",
							Description:             "Start of firewall rule range",
							CustomPropertyValidator: ipValidator,
						},
						"endIPAddress": &service.StringPropertySchema{
							Title:                   "End IP address",
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

func isGeoRedundentBackup(pp service.ProvisioningParameters) bool {
	return pp.GetString("backupRedundancy") == "geo"
}

func isSSLRequired(pp service.ProvisioningParameters) bool {
	return pp.GetString("sslEnforcement") != disabledParamString
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
