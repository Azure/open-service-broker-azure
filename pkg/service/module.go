package service

// Module is an interface to be implemented by the broker's modules
type Module interface {
	// GetName returns a module's name
	GetName() string
	// GetStability returns a module's relative level of stability
	GetStability() Stability
	// GetCatalog returns a Catalog of service/plans offered by a module
	GetCatalog() (Catalog, error)
}
