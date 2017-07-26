package echo

import (
	"context"
	"time"

	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
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
	params := provisioningParameters.(*ProvisioningParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) GetProvisioner() (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("generateMessageId", m.generateProvisioningMessageID),
		service.NewProvisioningStep("pause", m.pauseProvisioning),
		service.NewProvisioningStep("logMessage", m.logProvisioningMessage),
	)
}

func (m *module) generateProvisioningMessageID(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	pc := provisioningContext.(*ProvisioningContext)
	pc.MessageID = uuid.NewV4().String()
	return pc, nil
}

func (m *module) pauseProvisioning(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	select {
	case <-time.NewTimer(time.Minute).C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return provisioningContext, nil
}

func (m *module) logProvisioningMessage(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	pc := provisioningContext.(*ProvisioningContext)
	params := provisioningParameters.(*ProvisioningParameters)
	pc.Message = params.Message
	log.WithFields(log.Fields{
		"messageID": pc.MessageID,
		"message":   params.Message,
	}).Debug("provisioning instance")
	return pc, nil
}

func (m *module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	params := bindingParameters.(*BindingParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, error) {
	log.Info("Executing Bind...")
	params := bindingParameters.(*BindingParameters)
	bc := &BindingContext{
		MessageID: uuid.NewV4().String(),
		Message:   params.Message,
	}
	log.WithFields(log.Fields{
		"messageID": bc.MessageID,
		"message":   params.Message,
	}).Debug("binding instance")
	return bc, nil
}

func (m *module) Unbind(
	provisioningContext interface{},
	bindingContext interface{},
) error {
	bc := bindingContext.(*BindingContext)
	log.WithField(
		"messageID",
		bc.MessageID,
	).Debug("unbinding instance")
	return nil
}

func (m *module) GetDeprovisioner() (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("pause", m.pauseDeprovisioning),
		service.NewDeprovisioningStep(
			"logDeprovisioningMessage",
			m.logDeprovisioningMessage,
		),
	)
}

func (m *module) pauseDeprovisioning(
	ctx context.Context,
	provisioningContext interface{},
) (interface{}, error) {
	select {
	case <-time.NewTimer(time.Minute).C:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return provisioningContext, nil
}

func (m *module) logDeprovisioningMessage(
	ctx context.Context,
	provisioningContext interface{},
) (interface{}, error) {
	pc := provisioningContext.(*ProvisioningContext)
	log.WithField(
		"messageID",
		pc.MessageID,
	).Debug("deprovisioning instance")
	return pc, nil
}
