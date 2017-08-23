package fake

// DeleteFunction describes a function used to provide pluggable delete behavior
// to the fake implementation of the postgresql.Manager interface
type DeleteFunction func(serverName string, resourceGroupName string) error

// Manager is a fake implementaton of postgresql.Manager used for testing
type Manager struct {
	DeleteBehavior DeleteFunction
}

// NewManager returns a new, fake implementation of postgresql.Manager used for
// testing
func NewManager() *Manager {
	return &Manager{
		DeleteBehavior: defaultDeleteBehavior,
	}
}

// Delete simulates deletion of a PostgreSQL server
func (m *Manager) Delete(
	serverName string,
	resourceGroupName string,
) error {
	return m.DeleteBehavior(serverName, resourceGroupName)
}

func defaultDeleteBehavior(
	serverName string,
	resourceGroupName string,
) error {
	return nil
}
