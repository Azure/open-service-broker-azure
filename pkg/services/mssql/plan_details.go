package mssql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const (
	defaultCores            = 2
	defaultVCoreStorageInGB = 5
	maxStorageInGB          = 1024

	gen5Hardware = "Gen5"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getTierProvisionParameters(
		service.Instance,
	) (map[string]interface{}, error)
}

type dtuPlanDetails struct {
	//sku        string
	tier        string
	skuMap      map[int64]string
	allowedDTUs []int64
	defaultDTUs int64
	storageInGB int64
	includeDBMS bool
}

func addDBMSParameters(schema map[string]service.PropertySchema) {
	dbmsSchema := getDBMSCommonProvisionParamSchema().PropertySchemas
	for key, value := range dbmsSchema {
		schema[key] = value
	}
}

func (d dtuPlanDetails) getProvisionSchema() service.InputParametersSchema {
	schema := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if len(d.allowedDTUs) > 0 { //basic doesn't have DTUs, so don't add if not set
		schema.PropertySchemas["dtu"] = service.IntPropertySchema{
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	if d.includeDBMS {
		addDBMSParameters(schema.PropertySchemas)
	}
	return schema
}

func (d dtuPlanDetails) getTierProvisionParameters(
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
	p["sku"] = d.getSKU(pp)
	p["tier"] = d.tier
	p["maxSizeBytes"] = convertBytesToGB(d.storageInGB)
	return p, nil
}

func (d dtuPlanDetails) getSKU(pp databaseProvisionParams) string {
	if pp.Cores != nil {
		return d.skuMap[*pp.Cores]
	}
	return d.skuMap[d.defaultDTUs]
}

type vCorePlanDetails struct {
	tier          string
	tierShortName string
	includeDBMS   bool
}

func (v vCorePlanDetails) getProvisionSchema() service.InputParametersSchema {
	schema := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"cores": service.IntPropertySchema{
				AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
				DefaultValue:  ptr.ToInt64(defaultCores),
				Description:   "A virtual core represents the logical CPU",
			},
			"storage": service.IntPropertySchema{
				MinValue:     ptr.ToInt64(defaultVCoreStorageInGB),
				MaxValue:     ptr.ToInt64(maxStorageInGB),
				DefaultValue: ptr.ToInt64(defaultVCoreStorageInGB),
				Description:  "The maximum data storage capacity (in GB)",
			},
		},
	}

	// Include the DBMS params here if the plan details call for it
	if v.includeDBMS {
		addDBMSParameters(schema.PropertySchemas)
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

func convertBytesToGB(gb int64) int64 {
	return gb * 1024 * 1024 * 1024
}

func getStorageInBytes(
	pp databaseProvisionParams,
) int64 {
	if pp.Storage != nil {
		storageGB := *pp.Storage
		return convertBytesToGB(storageGB)
	}
	return convertBytesToGB(defaultVCoreStorageInGB)

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
