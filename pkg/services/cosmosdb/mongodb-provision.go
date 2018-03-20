package cosmosdb

import (
	"context"
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
	dt := cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	outputs, err := m.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"name": dt.DatabaseAccountName,
			"kind": "MongoDB",
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}

	sdt := cosmosdbSecureInstanceDetails{
		PrimaryKey: primaryKey,
	}

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

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}
