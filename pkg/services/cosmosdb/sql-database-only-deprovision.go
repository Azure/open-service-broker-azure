package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlDatabaseManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteDatabase", s.deleteDatabase),
	)
}

func (s *sqlDatabaseManager) deleteDatabase(
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
	err = deleteDatabase(pdt.DatabaseAccountName, dt.DatabaseName, psdt.PrimaryKey)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
