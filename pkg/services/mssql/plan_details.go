package mssql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getProvisionParameters(service.Instance) (map[string]interface{}, error)
}

type legacyPlanDetails struct {
	sku        string
	tier       string
	maxStorage int64
}

func (l legacyPlanDetails) getProvisionSchema() service.InputParametersSchema {
	return getDBMSCommonProvisionParamSchema()
}

func (l legacyPlanDetails) getProvisionParameters(
	instance service.Instance,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = l.sku
	p["tier"] = l.tier
	p["maxSizeBytes"] = l.maxStorage
	return p, nil
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
