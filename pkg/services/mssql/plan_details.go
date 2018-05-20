package mssql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const (
	defaultCores          = 2
	defaultStorageInGB    = 5
	defaultStorageInBytes = 5368709120
	maxStorageInBytes     = 1099511627776
	maxStorageInGB        = 1024

	gen5Hardware = "Gen5"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getTierProvisionParameters(
		service.Instance,
	) (map[string]interface{}, error)
}

type legacyPlanDetails struct {
	sku        string
	tier       string
	maxStorage int64
}

func (l legacyPlanDetails) getProvisionSchema() service.InputParametersSchema {
	return getDBMSCommonProvisionParamSchema()
}

func (l legacyPlanDetails) getTierProvisionParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = l.sku
	p["tier"] = l.tier
	p["maxSizeBytes"] = l.maxStorage
	return p, nil
}

type vCorePlanDetails struct {
	tier          string
	tierShortName string
	includesDBMS  bool
}

func (v vCorePlanDetails) getProvisionSchema() service.InputParametersSchema {
	schema := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"cores": service.IntPropertySchema{
				AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
				DefaultValue:  ptr.ToInt64(defaultCores),
				Description:   "A virtual core represents the logical CPU",
			},
			"storage": service.FloatPropertySchema{
				MinValue:     ptr.ToFloat64(defaultStorageInGB),
				MaxValue:     ptr.ToFloat64(maxStorageInGB),
				DefaultValue: ptr.ToFloat64(defaultStorageInGB),
				Description:  "The maximum data storage capacity",
			},
			"hardwareFamily": service.StringPropertySchema{
				AllowedValues: []string{"gen4", "gen5"},
				DefaultValue:  "gen5",
				Description: "Specifies the compute generation to use for " +
					"new instance",
			},
		},
	}

	// Include the DBMS params here if the plan details call for it
	if v.includesDBMS {
		dbmsSchema := getDBMSCommonProvisionParamSchema().PropertySchemas
		for key, value := range dbmsSchema {
			schema.PropertySchemas[key] = value
		}
	}
	return schema
}

func (v vCorePlanDetails) getTierProvisionParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	pp := databaseProvisionParams{}
	if err := service.GetStructFromMap(
		instance.ProvisioningParameters,
		&pp,
	); err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["sku"] = v.getSKU(pp)
	p["tier"] = v.tier
	p["maxSizeBytes"] = getStorageInBytes(pp)
	return p, nil
}

func (v vCorePlanDetails) getSKU(pp databaseProvisionParams) string {
	return fmt.Sprintf(
		"%s_%s_%d",
		v.tierShortName,
		gen5Hardware,
		getCores(pp),
	)
}

func getCores(pp databaseProvisionParams) int64 {
	if pp.Cores != nil {
		return *pp.Cores
	}
	return defaultCores
}

func getStorageInBytes(pp databaseProvisionParams) int64 {
	if pp.Storage != nil {
		storageGB := *pp.Storage
		return storageGB * 1024 * 1024 * 1024
	}
	return defaultStorageInBytes

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

func getDBMSCommonProvisionParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
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
			"sslEnforcement": &service.StringPropertySchema{
				Description: "Specifies whether the server requires the use of TLS" +
					" when connecting. Left unspecified, SSL will be enforced",
				AllowedValues: []string{"enabled", "disabled"},
				DefaultValue:  "enabled",
			},
		},
	}
}
