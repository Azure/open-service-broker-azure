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
	return "echo"
}

func (m *module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	params := provisioningParameters.(*echoProvisioningParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep(
			"generateMessageId",
			m.generateProvisioningMessageID,
		),
		service.NewProvisioningStep("pause", m.pauseProvisioning),
		service.NewProvisioningStep("logMessage", m.logProvisioningMessage),
	)
}

func (m *module) generateProvisioningMessageID(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	pc := provisioningContext.(*echoProvisioningContext)
	pc.MessageID = uuid.NewV4().String()
	return pc, nil
}

func (m *module) pauseProvisioning(
	ctx context.Context,
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	select {
	case <-time.NewTimer(time.Minute).C:
	case <-ctx.Done():
		log.Debug("context canceled; absorting pause")
		return nil, ctx.Err()
	}
	return provisioningContext, nil
}

func (m *module) logProvisioningMessage(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	pc := provisioningContext.(*echoProvisioningContext)
	params := provisioningParameters.(*echoProvisioningParameters)
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
	params := bindingParameters.(*echoBindingParameters)
	if params.Message == "bad message" {
		return service.NewValidationError("message", "message is a bad message!")
	}
	return nil
}

func (m *module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error) {
	messageID := uuid.NewV4().String()
	params := bindingParameters.(*echoBindingParameters)
	bc := &echoBindingContext{
		MessageID: messageID,
		Message:   params.Message,
	}
	c := &echoCredentials{
		MessageID: messageID,
	}
	log.WithFields(log.Fields{
		"messageID": bc.MessageID,
		"message":   params.Message,
	}).Debug("binding instance")
	return bc, c, nil
}

func (m *module) Unbind(
	provisioningContext interface{}, // nolint: unparam
	bindingContext interface{},
) error {
	bc := bindingContext.(*echoBindingContext)
	log.WithField(
		"messageID",
		bc.MessageID,
	).Debug("unbinding instance")
	return nil
}

func (m *module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
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
		log.Debug("context canceled; absorting pause")
		return nil, ctx.Err()
	}
	return provisioningContext, nil
}

func (m *module) logDeprovisioningMessage(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
) (interface{}, error) {
	pc := provisioningContext.(*echoProvisioningContext)
	log.WithField(
		"messageID",
		pc.MessageID,
	).Debug("deprovisioning instance")
	return pc, nil
}
