package storage

import "github.com/Azure/open-service-broker-azure/pkg/service"

// Store is an interface to be implemented by types capable of handling
// persistence for other broker-related types
type Store interface {
	// WriteInstance persists the given instance to the underlying storage
	WriteInstance(instance service.Instance) error
	// GetInstance retrieves a persisted instance from the underlying storage by
	// instance id
	GetInstance(instanceID string) (service.Instance, bool, error)
	// GetInstanceByID retrieves a persisted instance from the underlying storage
	// by alias
	GetInstanceByAlias(alias string) (service.Instance, bool, error)
	// GetInstanceChildCountByAlias returns the number of child instances
	GetInstanceChildCountByAlias(alias string) (int64, error)
	// DeleteInstance deletes a persisted instance from the underlying storage by
	// instance id
	DeleteInstance(instanceID string) (bool, error)
	// WriteBinding persists the given binding to the underlying storage
	WriteBinding(binding service.Binding) error
	// GetBinding retrieves a persisted instance from the underlying storage by
	// binding id
	GetBinding(bindingID string) (service.Binding, bool, error)
	// DeleteBinding deletes a persisted binding from the underlying storage by
	// binding id
	DeleteBinding(bindingID string) (bool, error)
	// TestConnection tests the connection to the underlying database (if there
	// is one)
	TestConnection() error
}
