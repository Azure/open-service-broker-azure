package service

const (
	// InstanceStateProvisioning represents the state where service instance
	// provisioning is in progress
	InstanceStateProvisioning = "PROVISIONING"
	// InstanceStateProvisioned represents the state where service instance
	// provisioning has completed successfully
	InstanceStateProvisioned = "PROVISIONED"
	// InstanceStateProvisioningFailed represents the state where service instance
	// provisioning has failed
	InstanceStateProvisioningFailed = "PROVISIONING_FAILED"
	// InstanceStateDeprovisioning represents the state where service instance
	// deprovisioning is in progress
	InstanceStateDeprovisioning = "DEPROVISIONING"
	// InstanceStateDeprovisioned represents the state where service instance
	// deprovisioning has completed successfully
	InstanceStateDeprovisioned = "DEPROVISIONED"
	// InstanceStateDeprovisioningFailed represents the state where service
	// instance deprovisioning has failed
	InstanceStateDeprovisioningFailed = "DEPROVISIONING_FAILED"
	// BindingStateBinding represents the state where service binding is in
	// progress
	BindingStateBinding = "BINDING"
	// BindingStateBound represents the state where service binding has completed
	// successfully
	BindingStateBound = "BOUND"
	// BindingStateBindingFailed represents the state where service binding has
	// failed
	BindingStateBindingFailed = "BINDING"
	// BindingStateUnbinding represents the state where service unbinding is in
	// progress
	BindingStateUnbinding = "UNBINDING"
	// BindingStateUnbound represents the state where service unbinding has
	// completed successfully
	BindingStateUnbound = "UNBOUND"
	// BindingStateUnbindingFailed represents the state where service unbinding
	// has failed
	BindingStateUnbindingFailed = "UNBINDING_FAILED"
)
