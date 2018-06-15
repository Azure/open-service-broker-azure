// +build experimental

package cosmosdb

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *sqlAllInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
		service.NewProvisioningStep("createDatabase", s.createDatabase),
	)
}

func (s *sqlAllInOneManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := sqlAllInOneInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(instance.Location),
		DatabaseName:        uuid.NewV4().String(),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *sqlAllInOneManager) deployARMTemplate(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {

	pp := &provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, pp); err != nil {
		return nil, nil, err
	}
	dt := &sqlAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	p, err := s.cosmosAccountManager.buildGoTemplateParams(
		instance,
		"GlobalDocumentDB",
	)
	if err != nil {
		return nil, nil, err
	}
	if instance.Tags == nil {
		instance.Tags = make(map[string]string)
	}

	fqdn, sdt, err := s.cosmosAccountManager.deployARMTemplate(ctx, instance, p)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	dt.FullyQualifiedDomainName = fqdn
	sdt.ConnectionString = fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
		dt.FullyQualifiedDomainName,
		sdt.PrimaryKey,
	)

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}

func (s *sqlAllInOneManager) createDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := &sqlAllInOneInstanceDetails{}
	err := service.GetStructFromMap(instance.Details, &dt)
	if err != nil {
		fmt.Printf("Failed to get DT Struct from Map %s", err)
		return nil, nil, err
	}
	sdt := &cosmosdbSecureInstanceDetails{}
	err = service.GetStructFromMap(instance.SecureDetails, &sdt)
	if err != nil {
		fmt.Printf("Failed to get SDT Struct from Map %s", err)
		return nil, nil, err
	}
	err = createDatabase(dt.DatabaseAccountName, dt.DatabaseName, sdt.PrimaryKey)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
