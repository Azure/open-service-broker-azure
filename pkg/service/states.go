package service

const (
	// InstanceStateProvisioningDeferred represents the state where service
	// instance provisioning has been requested and deferred pending the
	// completion of some other action
	InstanceStateProvisioningDeferred = "PROVISIONING_DEFERRED"
	// InstanceStateProvisioning represents the state where service instance
	// provisioning is in progress
	InstanceStateProvisioning = "PROVISIONING"
	// InstanceStateProvisioned represents the state where service instance
	// provisioning has completed successfully
	InstanceStateProvisioned = "PROVISIONED"
	// InstanceStateProvisioningFailed represents the state where service instance
	// provisioning has failed
	InstanceStateProvisioningFailed = "PROVISIONING_FAILED"
	// InstanceStateUpdating represents the state where service instance
	// updating is in progress
	InstanceStateUpdating = "UPDATING"
	// InstanceStateUpdated represents the state where service instance
	// updating has completed successfully
	// It redirects to InstanceStateProvisioned because it means the same thing
	// to any other operations besides updating
	InstanceStateUpdated = InstanceStateProvisioned
	// InstanceStateUpdatingFailed represents the state where service instance
	// updating has failed
	InstanceStateUpdatingFailed = "UPDATING_FAILED"
	// InstanceStateDeprovisioningDeferred represents the state where service
	// instance deprovisioning has been requested and deferred pending the
	// completion of some other action
	InstanceStateDeprovisioningDeferred = "DEPROVISIONING_DEFERRED"
	// InstanceStateDeprovisioning represents the state where service instance
	// deprovisioning is in progress
	InstanceStateDeprovisioning = "DEPROVISIONING"
	// InstanceStateDeprovisioningFailed represents the state where service
	// instance deprovisioning has failed
	InstanceStateDeprovisioningFailed = "DEPROVISIONING_FAILED"
	// BindingStateBound represents the state where service binding has completed
	// successfully
	BindingStateBound = "BOUND"
	// BindingStateBindingFailed represents the state where service binding has
	// failed
	BindingStateBindingFailed = "BINDING_FAILED"
	// BindingStateUnbindingFailed represents the state where service unbinding
	// has failed
	BindingStateUnbindingFailed = "UNBINDING_FAILED"
)
