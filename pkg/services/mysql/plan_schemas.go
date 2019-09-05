package mysql

import (
	"bytes"
	"fmt"
	"net"
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

// nolint: lll
func generateProvisioningParamsSchema(
	td tierDetails,
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
	ips.PropertySchemas["serverName"] = &service.StringPropertySchema{
		Title:       "Server Name",
		Description: "Name of the MySQL server",
		MinLength:   ptr.ToInt(3),
		MaxLength:   ptr.ToInt(63),
		// The server name can only contain lower case characters and numbers.
		AllowedPattern: `^[a-z0-9]+$`,
	}
	ips.PropertySchemas["adminAccountSettings"] = &service.ObjectPropertySchema{
		Title:       "Admin Account Setttings",
		Description: "Settings of administrator account of MySQL server. Typically you do not need to specify this.",
		PropertySchemas: map[string]service.PropertySchema{
			"adminUsername": &service.StringPropertySchema{
				Title:                   "Admin Username",
				Description:             "The administrator username for the server.",
				MinLength:               ptr.ToInt(1),
				MaxLength:               ptr.ToInt(63),
				CustomPropertyValidator: usernameValidator,
			},
			"adminPassword": &service.StringPropertySchema{
				Title:                   "Admin Password",
				Description:             "The administrator password for the server. **Warning**: you may leak your password if you specify this property, others can see this password in your request body and `ServiceInstance` definition. DO NOT use this property unless you know what you are doing.",
				MinLength:               ptr.ToInt(8),
				MaxLength:               ptr.ToInt(128),
				CustomPropertyValidator: passwordValidator,
			},
		},
	}
	ips.PropertySchemas["tags"] = &service.ObjectPropertySchema{
		Title: "Tags",
		Description: "Tags to be applied to new resources," +
			" specified as key/value pairs.",
		Additional: &service.StringPropertySchema{},
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
				DefaultValue: schemas.EnabledParamString,
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

// nolint: lll
func getBindingParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"username": &service.StringPropertySchema{
				Title:                   "Username",
				Description:             "The username to access created database.",
				MinLength:               ptr.ToInt(1),
				MaxLength:               ptr.ToInt(63),
				CustomPropertyValidator: usernameValidator,
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

// usernameValidator validates MySQL username,
// the username should:
// 1. It's a SQL Identifier, and not a typical system name
// (like admin, administrator, sa, root, dbmanager, loginmanager, etc.)
// 2. It shouldn't be a built-in database user or role
// (like dbo, guest, public, etc.)
// 3. It shouldn't contain whitespaces, unicode characters,
// or nonalphabetic characters,
// and that it doesn't begin with numbers or symbols.
func usernameValidator(context, value string) error {
	if value == "admin" ||
		value == "administrator" ||
		value == "sa" ||
		value == "root" ||
		value == "dbmanager" ||
		value == "loginmanager" ||
		value == "dbo" ||
		value == "guest" ||
		value == "public" {
		return service.NewValidationError(
			context,
			fmt.Sprintf("admin username can't be %s", value),
		)
	}
	runes := []rune(value)
	if !unicode.IsLetter(runes[0]) {
		return service.NewValidationError(
			context,
			fmt.Sprintf("admin username must begin with a character"),
		)
	}
	// For constraint3, it's really ambiguious which character is allowed.
	// When tested on Azure Portal, character `!@$` is allowed,
	// but character `#` is not allowed. I directly skip validating here
	// and let MySQL RP return the error.
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
// nolint: lll
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
	// Note: here if we don't add nolint, linter
	// will report error "should range over string, not []rune(string)"
	// The difference can be found here:
	// https://stackoverflow.com/questions/49062100/is-there-any-difference-between-range-str-and-range-runestr-in-golang
	// In our senario, that's OK to range over rune slice.
	for _, r := range runes { // nolint: megacheck
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
