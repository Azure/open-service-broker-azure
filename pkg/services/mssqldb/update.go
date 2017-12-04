package mssqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (s *serviceManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("update", s.update),
	)
}

func (s *serviceManager) update(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.UpdatingParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}

	if plan.GetID() == "" {
		return pc, nil
	}

	// Update service plan
	var resourceGroup string
	var location string
	if pc.IsNewServer {
		resourceGroup = standardProvisioningContext.ResourceGroup
		location = standardProvisioningContext.Location
	} else {
		servers := s.mssqlConfig.Servers
		server, ok := servers[pc.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pc.ServerName,
			)
		}
		resourceGroup = server.ResourceGroupName
		location = server.Location
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	if _, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		resourceGroup,
		location,
		armTemplateExistingServerBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"serverName":   pc.ServerName,
			"databaseName": pc.DatabaseName,
			"edition":      plan.GetProperties().Extended["edition"],
			"requestedServiceObjectiveName": plan.GetProperties().
				Extended["requestedServiceObjectiveName"],
			"maxSizeBytes": plan.GetProperties().
				Extended["maxSizeBytes"],
		},
		standardProvisioningContext.Tags,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	return pc, nil
}
