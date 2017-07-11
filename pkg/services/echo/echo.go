package echo

import (
	"log"

	"time"

	"context"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

type module struct{}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() service.Module {
	return &module{}
}

func (m *module) GetName() string {
	return "example"
}

func (m *module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	params := provisioningParameters.(echoProvisioningParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) GetProvisioner() (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("generateMessageId", m.generateProvisioningMessageID),
		service.NewProvisioningStep("pause", m.pause),
		service.NewProvisioningStep("logMessage", m.logProvisioningMessage),
	)
}

func (m *module) generateProvisioningMessageID(
	ctx context.Context,
	provisioningResult interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing generateProvisioningMessageID...")
	result := provisioningResult.(*echoProvisioningResult)
	result.MessageID = uuid.NewV4().String()
	return result, nil
}

func (m *module) pause(
	ctx context.Context,
	provisioningResult interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing pause...")
	select {
	case <-time.NewTimer(time.Minute).C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return provisioningResult, nil
}

func (m *module) logProvisioningMessage(
	ctx context.Context,
	provisioningResult interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	log.Println("Executing logProvisioningMessage...")
	result := provisioningResult.(*echoProvisioningResult)
	params := provisioningParameters.(*echoProvisioningParameters)
	result.Message = params.Message
	log.Printf("Provisioning %s: %s", result.MessageID, params.Message)
	return result, nil
}

func (m *module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	params := bindingParameters.(*echoBindingParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) Bind(
	provisioningResult interface{},
	bindingParameters interface{},
) (interface{}, error) {
	log.Println("Executing Bind...")
	params := bindingParameters.(*echoBindingParameters)
	result := &echoBindingResult{
		MessageID: uuid.NewV4().String(),
		Message:   params.Message,
	}
	log.Printf("Binding %s: %s", result.MessageID, params.Message)
	return result, nil
}

func (m *module) Unbind(
	provisioningResult interface{},
	bindingResult interface{},
) error {
	log.Println("Executing Unbind...")
	result := bindingResult.(*echoBindingResult)
	log.Printf("Unbinding %s: %s", result.MessageID, result.Message)
	return nil
}

func (m *module) GetDeprovisioner() (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep(
			"logDeprovisioningMessage",
			m.logDeprovisioningMessage,
		),
	)
}

func (m *module) logDeprovisioningMessage(
	ctx context.Context,
	provisioningResult interface{},
) (interface{}, error) {
	log.Println("Executing logDeprovisioningMessage...")
	result := provisioningResult.(*echoProvisioningResult)
	log.Printf("Deprovisioning %s: %s", result.MessageID, result.Message)
	return result, nil
}
