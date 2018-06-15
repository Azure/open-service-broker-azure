package service

// Module is an interface to be implemented by the broker's modules
type Module interface {
	// GetName returns a module's name
	GetName() string
	// GetCatalog returns a Catalog of service/plans offered by a module
	GetCatalog() (Catalog, error)
}
