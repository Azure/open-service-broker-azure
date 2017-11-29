package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	// No validation needed
	return nil
}

func (m *module) GetProvisioner(
	serviceID string,
	planID string,
) (service.Provisioner, error) {
	provisioningSteps := []service.ProvisioningStep{
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	}

	plan, err := m.getPlan(serviceID, planID)
	if err != nil {
		return nil, errors.New(
			"error getting plan by service ID and plan ID",
		)
	}
	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving the storage kind from the plan: %s",
			err,
		)
	}

	// Add provisioning steps that are specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		provisioningSteps = append(
			provisioningSteps,
			service.NewProvisioningStep("createBlobContainer", m.createBlobContainer),
		)
	}

	return service.NewProvisioner(provisioningSteps...)
}

func (m *module) preProvision(
	_ context.Context,
	_ string, // instanceID
	serviceID string,
	planID string,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.StorageAccountName = generate.NewIdentifier()

	plan, err := m.getPlan(serviceID, planID)
	if err != nil {
		return nil, errors.New(
			"error getting plan by service ID and plan ID",
		)
	}
	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving the storage kind from the plan: %s",
			err,
		)
	}

	// Add context that is specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		pc.ContainerName = uuid.NewV4().String()
	}

	return pc, nil
}

func (m *module) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	serviceID string,
	planID string,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}
	plan, err := m.getPlan(serviceID, planID)
	if err != nil {
		return nil, errors.New(
			"error getting plan by service ID and plan ID",
		)
	}
	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving the storage kind from the plan: %s",
			err,
		)
	}

	var armTemplateBytes []byte
	switch storeKind {
	case storageKindGeneralPurposeStorageAcccount:
		armTemplateBytes = armTemplateBytesGeneralPurposeStorage
	case storageKindBlobStorageAccount, storageKindBlobContainer:
		armTemplateBytes = armTemplateBytesBlobStorage
	}
	armTemplateParameters := map[string]interface{}{
		"name": pc.StorageAccountName,
	}
	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		armTemplateParameters, // ARM template params
		standardProvisioningContext.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	pc.AccessKey, ok = outputs["accessKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary access key from deployment: %s",
			err,
		)
	}

	return pc, nil
}

func (m *module) createBlobContainer(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}

	client, _ := storage.NewBasicClient(pc.StorageAccountName, pc.AccessKey)
	blobCli := client.GetBlobService()
	container := blobCli.GetContainerReference(pc.ContainerName)
	options := storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypePrivate,
	}
	_, err := container.CreateIfNotExists(&options)
	if err != nil {
		return nil, errors.New(
			"error creating container",
		)
	}

	return pc, nil
}

func (m *module) getPlan(serviceID, planID string) (service.Plan, error) {
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}

	service, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}

	plan, ok := service.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}

	return plan, nil
}
