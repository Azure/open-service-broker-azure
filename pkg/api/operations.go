package api

const (
	// OperationProvisioning represents the "provisioning" operation
	OperationProvisioning = "provisioning"
	// OperationUpdating represents the "updating" operation
	OperationUpdating = "updating"
	// OperationDeprovisioning represents the "deprovisioning" operation
	OperationDeprovisioning = "deprovisioning"
	// OperationStateInProgress represents the state of an operation that is still
	// pending completion
	OperationStateInProgress = "in progress"
	// OperationStateSucceeded represents the state of an operation that has
	// completed successfully
	OperationStateSucceeded = "succeeded"
	// OperationStateFailed represents the state of an operation that has
	// failed
	OperationStateFailed = "failed"
	// OperationStateGone is a pseudo oepration state represting the "state"
	// of an operation against an entity that no longer exists
	OperationStateGone = "gone"
)
