package postgresql

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"unicode"

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
	allowedBackupRedundancy []service.EnumValue
}

func (t *tierDetails) getSku(pp service.ProvisioningParameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.

	// Temporary workaround for Mooncake
	var skuFamily string
	location := pp.GetString("location")
	if location == "chinanorth" || location == "chinaeast" {
		skuFamily = "Gen4"
	} else {
		skuFamily = "Gen5"
	}

	sku := fmt.Sprintf(
		"%s_%s_%d",
		t.tierShortName,
		skuFamily,
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
		Title:        "Backup redundancy",
		Description:  "Specifies the backup redundancy",
		OneOf:        td.allowedBackupRedundancy,
		DefaultValue: "local",
	}
	ips.PropertySchemas["tags"] = &service.ObjectPropertySchema{
		Title: "Tags",
		Description: "Tags to be applied to new resources," +
			" specified as key/value pairs.",
		Additional: &service.StringPropertySchema{},
	}
	ips.PropertySchemas["serverName"] = &service.StringPropertySchema{
		Title:          "Server Name",
		Description:    "Name of the postgreSQL server",
		MinLength:      ptr.ToInt(3),
		MaxLength:      ptr.ToInt(63),
		AllowedPattern: `^[a-z0-9]+[-a-z0-9]*[a-z0-9]+$`,
	}
	ips.PropertySchemas["adminAccountSettings"] = &service.ObjectPropertySchema{
		Title:       "Admin Account Setttings",
		Description: "Settings of administrator account of PostgreSQL server. Typically you do not need to specify this.",
		PropertySchemas: map[string]service.PropertySchema{
			"adminUsername": &service.StringPropertySchema{
				Title:          "Admin Username",
				Description:    "The administrator username for the server.",
				MinLength:      ptr.ToInt(1),
				MaxLength:      ptr.ToInt(63),
				AllowedPattern: `^(?!azure_superuser$)(?!azure_pg_admin$)(?!admin$)(?!administrator$)(?!root$)(?!guest$)(?!public$)(?!pg_)[_a-z0-9]+`,
			},
			"adminPassword": &service.StringPropertySchema{
				Title:                   "Admin Password",
				Description:             "The administrator password for the server. **Warning**: you may leak your password if you specify this property, others can see this password in your request body and `ServiceInstance` definition. DO NOT use this property unless you know what you are doing.",
				MinLength:               ptr.ToInt(8),
				MaxLength:               ptr.ToInt(128),
				CustomPropertyValidator: passwordValidator,
			},
		},
		CustomPropertyValidator: adminAccountSettingValidator,
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
				OneOf:        schemas.EnabledDisabledValues(),
				DefaultValue: schemas.DisabledParamString,
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
			"virtualNetworkRules": &service.ArrayPropertySchema{
				Title:       "Virtual network rules",
				Description: "Virtual network rules to apply to instance. ",
				ItemsSchema: &service.ObjectPropertySchema{
					Title:       "Virtual network rule",
					Description: "Individual virtual network rule",
					RequiredProperties: []string{
						"name",
						"subnetId",
					},
					PropertySchemas: map[string]service.PropertySchema{
						"name": &service.StringPropertySchema{
							Title:       "Name",
							Description: "Name of virtual network rule",
						},
						"subnetId": &service.StringPropertySchema{
							Title:       "Subnet ID",
							Description: "Subnet ID to add",
						},
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
	return pp.GetString("sslEnforcement") != schemas.DisabledParamString
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

// passwordValidator validates postgreSQL password,
// the password should:
// 1. Have at least 8 characters and at most 128 characters.
// 2. Contain characters from three of the following categories:
//    – English uppercase letters
//    - English lowercase letters
//    - numbers (0-9),
//    - non-alphanumeric characters (!, $, #, %, etc.).
func passwordValidator(context, value string) error {
	if len(value) < 8 || len(value) > 128 {
		return service.NewValidationError(
			context,
			fmt.Sprintf("the passsword should have at least 8 characters and at most 128 characters, given password's length is %d", len(value)), // nolint: lll
		)
	}
	runes := []rune(value)
	var (
		digitOccurred              int
		lowercaseCharacterOccurred int
		uppercaseCharacterOccurred int
		specialCharacterOccurred   int
	)
	for _, r := range runes {
		if unicode.IsDigit(r) {
			digitOccurred = 1
		} else if unicode.IsLower(r) {
			lowercaseCharacterOccurred = 1
		} else if unicode.IsUpper(r) {
			uppercaseCharacterOccurred = 1
		} else {
			specialCharacterOccurred = 1
		}
	}

	if digitOccurred+lowercaseCharacterOccurred+uppercaseCharacterOccurred+specialCharacterOccurred < 3 {
		return service.NewValidationError(
			context,
			"the password must contain characters from three of the following categories – English uppercase letters, English lowercase letters, numbers (0-9), and non-alphanumeric characters.", // nolint: lll
		)
	}
	return nil
}

// adminAccountSettingValidator, in fact, validates
// the admin password does not contain all or part of
// the admin username. Part of a login name is defined as
// three or more consecutive alphanumeric characters.
// The reason we need this function is we can't get the
// admin username in password's custom validator, so we
// have to wrap them into an object and validate it
// in this addtional validator.
func adminAccountSettingValidator(
	context string,
	valMap map[string]interface{},
) error {
	var username, password string

	usernameInterface := valMap["adminUsername"]
	passwordInterface := valMap["password"]
	// If user does not specify password, OSBA will
	// generate one for user, it has very little
	// possibility to conflict, we directly return nil here.
	if passwordInterface == nil {
		return nil
	}

	if usernameInterface == nil {
		username = "postgres"
	} else {
		username = usernameInterface.(string)
	}
	password = passwordInterface.(string)

	// Find whether password contains part of the username.
	// We only need to detect whether password contains username's
	// substring of length 3.
	// That's OK to interate over all username's substrings
	// of length 3 here as the stirng is really short. Though
	// AC automation can save some time, but code complexity
	// is too high, and I think we do not really need it here.
	found := false
	containedSubstr := ""
	usernameLen := len(username)
	for startIdx := 0; startIdx <= usernameLen-3; startIdx++ {
		subStr := username[startIdx : startIdx+3]
		found = strings.Contains(password, subStr)
		if found {
			containedSubstr = subStr
			break
		}
	}
	if found {
		return service.NewValidationError(
			context,
			fmt.Sprintf("the password should not contain part of username, username is %s, contained part is %s", username, containedSubstr), // nolint: lll
		)
	}
	return nil
}
