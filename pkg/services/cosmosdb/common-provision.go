package cosmosdb

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

var (
	validValues = map[string]struct{}{
		"eventual":         {},
		"session":          {},
		"boundedstaleness": {},
		"strong":           {},
		"consistentprefix": {},
	}
)

const (
	disabled = "disabled"
	enabled  = "enabled"
)

func generateAccountName(location string) string {
	databaseAccountName := uuid.NewV4().String()
	// CosmosDB currently limits database account names to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and generate a shorter random identifier if needed.
	effectiveNameLength := len(location) + len(databaseAccountName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		databaseAccountName = generate.NewIdentifierOfLength(nameLength)
		logFields := log.Fields{
			"name":   databaseAccountName,
			"length": len(databaseAccountName),
		}
		log.WithFields(logFields).Debug(
			"returning fallback database account name",
		)
	}
	return databaseAccountName
}

func preProvision(
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := cosmosdbInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(instance.Location),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (c *cosmosAccountManager) ValidateProvisioningParameters(
	provisionParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp := provisioningParameters{}
	if err := service.GetStructFromMap(provisionParameters, &pp); err != nil {
		return err
	}
	if pp.IPFilterRules != nil {
		allowAzure := strings.ToLower(pp.IPFilterRules.AllowAzure)
		if allowAzure != "" && allowAzure != enabled &&
			allowAzure != disabled {
			return service.NewValidationError(
				"allowAzure",
				fmt.Sprintf(`invalid option: "%s"`, pp.IPFilterRules.AllowAzure),
			)
		}
		allowPortal := strings.ToLower(pp.IPFilterRules.AllowPortal)
		if allowPortal != "" && allowPortal != enabled &&
			allowPortal != disabled {
			return service.NewValidationError(
				"allowPortal",
				fmt.Sprintf(`invalid option: "%s"`, pp.IPFilterRules.AllowPortal),
			)
		}
		for _, filter := range pp.IPFilterRules.Filters {
			// First check if it is a valid IP
			ip := net.ParseIP(filter)
			if ip == nil {
				// Check to see if it is a valid CIDR
				ip, _, _ = net.ParseCIDR(filter)
				if ip == nil {
					return service.NewValidationError(
						"IP Filter",
						fmt.Sprintf(`invalid value: "%s"`, filter),
					)
				}
			}
		}
	}
	if pp.ConsistencyPolicy != nil {
		cp := strings.ToLower(pp.ConsistencyPolicy.DefaultConsistency)
		if _, valid := validValues[cp]; !valid {
			return service.NewValidationError(
				"consistencyPolicy",
				fmt.Sprintf(`invalid default consistency level: "%s"`, cp),
			)
		}
		if cp == "boundedstaleness" {
			if pp.ConsistencyPolicy.BoundedStaleness == nil {
			} else {
				if pp.ConsistencyPolicy.BoundedStaleness.MaxInternal == nil {
					return service.NewValidationError(
						"consistencyPolicy",
						"maxIntervalInSeconds must be provided when "+
							"defaultConsistencyPolicy is set to "+
							"'BoundedStaleness'.",
					)
				}
				if pp.ConsistencyPolicy.BoundedStaleness.MaxStaleness == nil {
					return service.NewValidationError(
						"consistencyPolicy",
						"maxStalenessPrefix must be provided when "+
							"defaultConsistencyPolicy is set to "+
							"'BoundedStaleness'.",
					)
				}

				maxIntervalParam :=
					*pp.ConsistencyPolicy.BoundedStaleness.MaxInternal
				if maxIntervalParam < 5 || maxIntervalParam > 86400 {
					return service.NewValidationError(
						"consistencyPolicy",
						fmt.Sprintf(
							`invalid maxIntervalInSeconds: "%d"`,
							maxIntervalParam,
						),
					)
				}

				maxStalenessParam :=
					*pp.ConsistencyPolicy.BoundedStaleness.MaxStaleness
				if maxStalenessParam < 1 || maxStalenessParam > 2147483647 {
					return service.NewValidationError(
						"consistencyPolicy",
						fmt.Sprintf(
							`invalid maxStalenessPrefix: "%d"`,
							maxStalenessParam,
						),
					)
				}

			}
		}
	}
	return nil

}

func (c *cosmosAccountManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return preProvision(instance)
}

func (c *cosmosAccountManager) buildGoTemplateParams(
	instance service.Instance,
	kind string,
) (map[string]interface{}, error) {

	pp := &provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, pp); err != nil {
		return nil, err
	}

	dt := &sqlAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}

	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind

	filters := []string{}

	if pp.IPFilterRules != nil {
		allowAzure := strings.ToLower(pp.IPFilterRules.AllowPortal)
		allowPortal := strings.ToLower(pp.IPFilterRules.AllowPortal)
		if allowAzure != disabled {
			filters = append(filters, "0.0.0.0")
		} else if allowPortal != disabled {
			// Azure Portal IP Addresses per:
			// https://aka.ms/Vwxndo
			//|| Region            || IP address(es) ||
			//||=====================================||
			//|| China             || 139.217.8.252  ||
			//||===================||================||
			//|| Germany           || 51.4.229.218   ||
			//||===================||================||
			//|| US Gov            || 52.244.48.71   ||
			//||===================||================||
			//|| All other regions || 104.42.195.92  ||
			//||                   || 40.76.54.131   ||
			//||                   || 52.176.6.30    ||
			//||                   || 52.169.50.45   ||
			//||                   || 52.187.184.26  ||
			//=======================================||
			// Given that we don't really have context of the cloud
			// we are provisioning with right now, use all of the above
			// addresses.
			filters = append(filters,
				"104.42.195.92",
				"40.76.54.131",
				"52.176.6.30",
				"52.169.50.45",
				"52.187.184.26",
				"51.4.229.218",
				"139.217.8.252",
				"52.244.48.71",
			)
		}
		filters = append(filters, pp.IPFilterRules.Filters...)
	} else {
		filters = append(filters, "0.0.0.0")
	}
	if len(filters) > 0 {
		p["ipFilters"] = strings.Join(filters, ",")
	}

	if pp.ConsistencyPolicy != nil {
		consistencyPolicy := make(map[string]interface{})
		cp := strings.ToLower(pp.ConsistencyPolicy.DefaultConsistency)
		switch cp {
		case "eventual":
			consistencyPolicy["defaultConsistencyLevel"] = "Eventual"
		case "session":
			consistencyPolicy["defaultConsistencyLevel"] = "Session"
		case "boundedstaleness":
			consistencyPolicy["defaultConsistencyLevel"] = "BoundedStaleness"
		case "strong":
			consistencyPolicy["defaultConsistencyLevel"] = "Strong"
		case "consistentprefix":
			consistencyPolicy["defaultConsistencyLevel"] = "ConsistentPrefix"
		}

		if consistencyPolicy["defaultConsistencyLevel"] == "BoundedStaleness" {
			boundedStalenessSettings := make(map[string]interface{})
			boundedStalenessSettings["maxIntervalInSeconds"] =
				*pp.ConsistencyPolicy.BoundedStaleness.MaxInternal

			boundedStalenessSettings["maxStalenessPrefix"] =
				*pp.ConsistencyPolicy.BoundedStaleness.MaxStaleness
		}
		p["consistencyPolicy"] = consistencyPolicy

	}
	return p, nil
}

func (c *cosmosAccountManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	goParams map[string]interface{},
) (string, *cosmosdbSecureInstanceDetails, error) {
	dt := &cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return "", nil, err
	}
	fqdn, sdt, err := c.deployTemplate(
		instance,
		goParams,
		dt.ARMDeploymentName,
		armTemplateBytes,
	)
	if err != nil {
		return "", nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return fqdn, sdt, nil
}

func (c *cosmosAccountManager) deployTemplate(
	instance service.Instance,
	goParams map[string]interface{},
	armDeploymentName string,
	armTemplateBytes []byte,
) (string, *cosmosdbSecureInstanceDetails, error) {
	outputs, err := c.armDeployer.Deploy(
		armDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return "", nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return c.handleOutput(outputs)
}

func (c *cosmosAccountManager) handleOutput(
	outputs map[string]interface{},
) (string, *cosmosdbSecureInstanceDetails, error) {

	var ok bool
	fqdn, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return "", nil, fmt.Errorf("error retrieving fully qualified domain name from deployment")
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return "", nil, fmt.Errorf("error retrieving primary key from deployment")
	}

	sdt := cosmosdbSecureInstanceDetails{
		PrimaryKey: primaryKey,
	}

	return fqdn, &sdt, nil
}
