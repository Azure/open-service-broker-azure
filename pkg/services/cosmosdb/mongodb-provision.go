package cosmosdb

import (
	"context"
	"fmt"

	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	)
}

func (m *mongoAccountManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {

	pp := instance.ProvisioningParameters
	dt := instance.Details.(*cosmosdbInstanceDetails)
	p, err := m.buildGoTemplateParams(
		pp,
		dt,
		"MongoDB",
	)
	if err != nil {
		return nil, err
	}
	tags := getTags(pp)
	fqdn, pk, err := m.cosmosAccountManager.deployARMTemplate(
		pp,
		dt,
		p,
		tags,
	)

	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	dt.FullyQualifiedDomainName = fqdn
	dt.PrimaryKey = service.SecureString(pk)
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
	dt.ConnectionString = service.SecureString(
		fmt.Sprintf(
			"mongodb://%s:%s@%s:10255/?ssl=true&replicaSet=globaldb",
			dt.DatabaseAccountName,
			dt.PrimaryKey,
			dt.FullyQualifiedDomainName,
		),
	)

	return dt, err
}
