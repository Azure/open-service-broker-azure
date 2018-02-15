package cosmosdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseAccountName = generateDatabaseName(instance.Location)
	return dt, instance.SecureDetails, nil
}

func generateDatabaseName(location string) string {
	databaseName := uuid.NewV4().String()
	// CosmosDB currently limits database name to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and generate a shorter random identifier if needed.
	effectiveNameLength := len(location) + len(databaseName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		databaseName = generate.NewIdentifierOfLength(nameLength)
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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*cosmosdbSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as *cosmosdbSecureInstanceDetails",
		)
	}
	plan := instance.Plan
	dt.DatabaseKind, ok = plan.GetProperties().Extended[kindKey].(databaseKind)
	if !ok {
		return nil, nil, errors.New(
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
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	fullyQualifiedDomainName, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	dt.FullyQualifiedDomainName = fullyQualifiedDomainName

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	sdt.PrimaryKey = primaryKey

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
		sdt.ConnectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
			dt.DatabaseAccountName,
			sdt.PrimaryKey,
			dt.FullyQualifiedDomainName,
		)
	case databaseKindGlobalDocumentDB:
		sdt.ConnectionString = fmt.Sprintf(
			"AccountEndpoint=%s;AccountKey=%s;",
			dt.FullyQualifiedDomainName,
			sdt.PrimaryKey,
		)
	}

	return dt, sdt, nil
}
