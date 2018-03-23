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
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {

	pp := &provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, pp); err != nil {
		return nil, nil, err
	}

	dt := &cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, dt); err != nil {
		return nil, nil, err
	}

	p := m.buildGoTemplateParams(pp, dt)
	dt, sdt, err := m.cosmosAccountManager.deployARMTemplate(ctx, instance, p)

	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
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
