package cosmosdb

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

const ascii = "1234567890abcdefghijklmnopqrstuvwxyz"

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	// Nothing to validate
	return nil
}

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseAccountName = generateDatabaseName(instance.Location)
	return dt, nil
}

func generateDatabaseName(location string) string {
	databaseName := uuid.NewV4().String()
	// CosmosDB currently limits database name to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and truncate.
	effectiveNameLength := len(location) + len(databaseName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		b := make([]byte, nameLength)
		for i := range b {
			b[i] = ascii[rand.Intn(len(ascii))]
		}
		databaseName = string(b)
		logFields := log.Fields{
			"name":   databaseName,
			"length": len(databaseName),
		}
		log.WithFields(logFields).Debug("returning fallback database name")
	}
	return databaseName
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	dt.DatabaseKind, ok = plan.GetProperties().Extended[kindKey].(databaseKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the kind from deployment",
		)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"name": dt.DatabaseAccountName,
			"kind": plan.GetProperties().Extended[kindKey],
		},
		instance.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	fullyQualifiedDomainName, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	dt.FullyQualifiedDomainName = fullyQualifiedDomainName

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	dt.PrimaryKey = primaryKey

	switch dt.DatabaseKind {
	case databaseKindMongoDB:
		// Allow to remove the https:// and the port 443 on the FQDN
		// This will allow to adapt the FQDN for Azure Public / Azure Gov ...
		// Before :
		// https://6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com:443
		// After :
		// 6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com
		hostnameNoHTTPS := strings.Join(
			strings.Split(dt.FullyQualifiedDomainName, "https://"),
			"",
		)
		dt.FullyQualifiedDomainName = strings.Join(
			strings.Split(hostnameNoHTTPS, ":443/"),
			"",
		)
		dt.ConnectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
			dt.DatabaseAccountName,
			dt.PrimaryKey,
			dt.FullyQualifiedDomainName,
		)
	case databaseKindGlobalDocumentDB:
		dt.ConnectionString = fmt.Sprintf(
			"AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			dt.PrimaryKey,
		)
	}

	return dt, nil
}
