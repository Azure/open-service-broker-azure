package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *sqlDatabaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("createDatabase", s.createDatabase),
	)
}

func (s *sqlDatabaseManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := sqlAllInOneInstanceDetails{
		DatabaseName: uuid.NewV4().String(),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *sqlDatabaseManager) createDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := &sqlDatabaseOnlyInstanceDetails{}
	err := service.GetStructFromMap(instance.Details, &dt)
	if err != nil {
		fmt.Printf("Failed to get DT Struct from Map %s", err)
		return nil, nil, err
	}

	pdt := &cosmosdbInstanceDetails{}
	err = service.GetStructFromMap(instance.Parent.Details, &pdt)
	if err != nil {
		fmt.Printf("Failed to get PDT Struct from Map %s", err)
		return nil, nil, err
	}

	psdt := &cosmosdbSecureInstanceDetails{}
	err = service.GetStructFromMap(instance.Parent.SecureDetails, &psdt)
	if err != nil {
		fmt.Printf("Failed to get SDT Struct from Map %s", err)
		return nil, nil, err
	}
	err = createDatabase(pdt.DatabaseAccountName, dt.DatabaseName, psdt.PrimaryKey)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
