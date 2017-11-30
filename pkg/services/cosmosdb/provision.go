package cosmosdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
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
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *cosmosdbProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.DatabaseAccountName = uuid.NewV4().String()
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *cosmosdbProvisioningContext",
		)
	}
	pc.DatabaseKind, ok = plan.GetProperties().Extended[kindKey].(databaseKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the kind from deployment",
		)
	}

	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"name": pc.DatabaseAccountName,
			"kind": plan.GetProperties().Extended[kindKey],
		},
		standardProvisioningContext.Tags,
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
	pc.FullyQualifiedDomainName = fullyQualifiedDomainName

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	pc.PrimaryKey = primaryKey

	switch pc.DatabaseKind {
	case databaseKindMongoDB:
		// Allow to remove the https:// and the port 443 on the FQDN
		// This will allow to adapt the FQDN for Azure Public / Azure Gov ...
		// Before :
		// https://6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com:443
		// After :
		// 6bd965fd-a916-4c3c-9606-161ec4d726bf.documents.azure.com
		hostnameNoHTTPS := strings.Join(
			strings.Split(pc.FullyQualifiedDomainName, "https://"),
			"",
		)
		pc.FullyQualifiedDomainName = strings.Join(
			strings.Split(hostnameNoHTTPS, ":443/"),
			"",
		)
		pc.ConnectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
			pc.DatabaseAccountName,
			pc.PrimaryKey,
			pc.FullyQualifiedDomainName,
		)
	case databaseKindGlobalDocumentDB:
		pc.ConnectionString = fmt.Sprintf(
			"AccountEndpoint=%s;AccountKey=%s;",
			pc.FullyQualifiedDomainName,
			pc.PrimaryKey,
		)
	}

	return pc, nil
}
