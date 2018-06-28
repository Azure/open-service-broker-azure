package cosmosdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

const disabled = "disabled"

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

func (c *cosmosAccountManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	l := instance.ProvisioningParameters.GetString("location")
	return &cosmosdbInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(l),
	}, nil
}

func (c *cosmosAccountManager) buildGoTemplateParams(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind
	p["location"] = pp.GetString("location")
	filters := []string{}
	ipFilters := pp.GetObject("ipFilters")
	if ipFilters.GetString("allowAzure") != disabled {
		filters = append(filters, "0.0.0.0")
	} else if ipFilters.GetString("allowPortal") != disabled {
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
	} else {
		filters = append(filters, "0.0.0.0")
	}
	filters = append(filters, ipFilters.GetStringArray("allowedIPRanges")...)
	if len(filters) > 0 {
		p["ipFilters"] = strings.Join(filters, ",")
	}
	p["consistencyPolicy"] = pp.GetObject("consistencyPolicy").Data
	return p, nil
}

func getTags(pp *service.ProvisioningParameters) map[string]string {
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	return tags
}

func (c *cosmosAccountManager) deployARMTemplate(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	goParams map[string]interface{},
	tags map[string]string,
) (string, string, error) {
	outputs, err := c.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return "", "", fmt.Errorf("error deploying ARM template: %s", err)
	}
	fqdn, primaryKey, err := c.handleOutput(outputs)
	if err != nil {
		return "", "", fmt.Errorf("error deploying ARM template: %s", err)
	}
	return fqdn, primaryKey, nil
}

func (c *cosmosAccountManager) handleOutput(
	outputs map[string]interface{},
) (string, string, error) {

	var ok bool
	fqdn, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return "", "", fmt.Errorf(
			"error retrieving fully qualified domain name from deployment",
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return "", "", fmt.Errorf("error retrieving primary key from deployment")
	}
	return fqdn, primaryKey, nil
}
