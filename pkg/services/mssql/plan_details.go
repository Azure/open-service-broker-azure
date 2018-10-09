package mssql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getProvisionSchemaForExistingIntance() service.InputParametersSchema
	getTierProvisionParameters(
		pp service.ProvisioningParameters,
	) (map[string]interface{}, error)
	getUpdateSchema() service.InputParametersSchema
	validateUpdateParameters(service.Instance) error
}

type dtuPlanDetails struct {
	tierName    string
	skuMap      map[int64]string
	allowedDTUs []int64
	defaultDTUs int64
	storageInGB int64
	includeDBMS bool
}

func (d dtuPlanDetails) validateUpdateParameters(service.Instance) error {
	return nil // no op
}

func (d dtuPlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if d.includeDBMS {
		ips = getDBMSCommonUpdateParamSchema()
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Title:         "DTUs",
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if d.includeDBMS {
		ips = getDBMSCommonProvisionParamSchema()
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Title:         "DTUs",
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getProvisionSchemaForExistingIntance() service.InputParametersSchema { // nolint: lll
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if d.includeDBMS {
		ips = getDBMSCommonProvisionParamSchema()
	}
	ips.PropertySchemas["database"] = &service.StringPropertySchema{
		Title:       "Database",
		Description: "The name of the existing database",
	}
	return ips
}

func (d dtuPlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = d.getSKU(pp)
	p["tier"] = d.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (d dtuPlanDetails) getSKU(pp service.ProvisioningParameters) string {
	// Basic tier is constrained to just 5 DTUs, if this is the basic tier, there
	// is no dtus param. We can infer this is the case if the tier details don't
	// tell us there's a choice.
	if len(d.allowedDTUs) == 0 {
		return d.skuMap[d.defaultDTUs]
	}
	return d.skuMap[pp.GetInt64("dtus")]
}

type vCorePlanDetails struct {
	tierName      string
	tierShortName string
	includeDBMS   bool
}

func (v vCorePlanDetails) validateUpdateParameters(
	instance service.Instance,
) error {
	return validateStorageUpdate(
		*instance.ProvisioningParameters,
		*instance.UpdatingParameters,
	)
}

func (v vCorePlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if v.includeDBMS {
		ips = getDBMSCommonUpdateParamSchema()
	}
	ips.PropertySchemas["cores"] = &service.IntPropertySchema{
		AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
		DefaultValue:  ptr.ToInt64(2),
		Title:         "Cores",
		Description:   "A virtual core represents the logical CPU",
	}
	ips.PropertySchemas["storage"] = &service.IntPropertySchema{
		MinValue:     ptr.ToInt64(5),
		MaxValue:     ptr.ToInt64(1024),
		DefaultValue: ptr.ToInt64(10),
		Title:        "Storage",
		Description:  "The maximum data storage capacity (in GB)",
	}
	return ips
}

func (v vCorePlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if v.includeDBMS {
		ips = getDBMSCommonProvisionParamSchema()
	}
	ips.PropertySchemas["cores"] = &service.IntPropertySchema{
		AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
		DefaultValue:  ptr.ToInt64(2),
		Title:         "Cores",
		Description:   "A virtual core represents the logical CPU",
	}
	ips.PropertySchemas["storage"] = &service.IntPropertySchema{
		MinValue:     ptr.ToInt64(5),
		MaxValue:     ptr.ToInt64(1024),
		DefaultValue: ptr.ToInt64(10),
		Title:        "Storage",
		Description:  "The maximum data storage capacity (in GB)",
	}
	return ips
}

func (v vCorePlanDetails) getProvisionSchemaForExistingIntance() service.InputParametersSchema { // nolint: lll
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	if v.includeDBMS {
		ips = getDBMSCommonProvisionParamSchema()
	}
	ips.PropertySchemas["database"] = &service.StringPropertySchema{
		Title:       "Database",
		Description: "The name of the existing database",
	}
	return ips
}

func (v vCorePlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = v.getSKU(pp)
	p["tier"] = v.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (v vCorePlanDetails) getSKU(pp service.ProvisioningParameters) string {
	return fmt.Sprintf(
		"%s_Gen5_%d",
		v.tierShortName,
		pp.GetInt64("cores"),
	)
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

func getDBMSCommonUpdateParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
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

func getDBMSCommonProvisionParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{"location", "resourceGroup"},
		PropertySchemas: map[string]service.PropertySchema{
			"location": &service.StringPropertySchema{
				Title: "Location",
				Description: "The Azure region in which to provision" +
					" applicable resources.",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"resourceGroup": &service.StringPropertySchema{
				Title: "Resource group",
				Description: "The (new or existing) resource group with which" +
					" to associate new resources.",
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
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}

func getDBMSRegisteredUpdateParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		SecureProperties: []string{
			"administratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"administratorLogin": &service.StringPropertySchema{
				Title: "Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the existing server",
			},
			"administratorLoginPassword": &service.StringPropertySchema{
				Title: "Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the existing server",
			},
		},
	}
}

func getDBMSRegisteredProvisionParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"resourceGroup",
			"location",
			"server",
			"administratorLogin",
			"administratorLoginPassword",
		},
		SecureProperties: []string{
			"administratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"resourceGroup": &service.StringPropertySchema{
				Title:       "Resource Group",
				Description: "Specifies the resource group of the existing server",
			},
			"location": &service.StringPropertySchema{
				Title:       "Location",
				Description: "Specifies the location of the existing server",
			},
			"server": &service.StringPropertySchema{
				Title:       "Server Name",
				Description: "Specifies the name of the existing server",
			},
			"administratorLogin": &service.StringPropertySchema{
				Title: "Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the existing server",
			},
			"administratorLoginPassword": &service.StringPropertySchema{
				Title: "Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the existing server",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}

func validateStorageUpdate(
	pp service.ProvisioningParameters,
	up service.ProvisioningParameters,
) error {
	existingStorage := pp.GetInt64("storage")
	newStorge := up.GetInt64("storage")
	if newStorge < existingStorage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(
				`invalid value: cannot reduce storage from %d to %d`,
				existingStorage,
				newStorge,
			),
		)
	}
	return nil
}
