package cosmosdb

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) ValidateProvisioningParameters(
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
) error {
	// Nothing to validate
	return nil
}

func (m *mongoAccountManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	)
}

func (m *mongoAccountManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return preProvision(instance)
}

func (m *mongoAccountManager) deployARMTemplate(
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
	dt.DatabaseKind = "MongoDB"
	if !ok {
		return nil, nil, errors.New(
			"error retrieving the kind from deployment",
		)
	}

	outputs, err := m.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"name": dt.DatabaseAccountName,
			"kind": dt.DatabaseKind,
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

	return dt, sdt, nil
}
