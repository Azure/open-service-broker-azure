package service

import (
	"encoding/json"
)

// Catalog is an interface to be implemented by types that represents the
// service/plans offered by a service module or by the entire broker.
type Catalog interface {
	GetServices() []Service
	GetService(serviceID string) (Service, bool)
}

type catalog struct {
	services        []Service
	indexedServices map[string]Service
}

// ServiceProperties represent the properties of a Service that can be directly
// instantiated and passed to the NewService() constructor function which will
// carry out all necessary initialization.
type ServiceProperties struct { // nolint: golint
	Name          string          `json:"name"`
	ID            string          `json:"id"`
	Description   string          `json:"description"`
	Metadata      ServiceMetadata `json:"metadata,omitempty"`
	Tags          []string        `json:"tags"`
	Bindable      bool            `json:"bindable"`
	PlanUpdatable bool            `json:"plan_updateable"` // Misspelling is
	// deliberate to match the spec
	ParentServiceID string                 `json:"-"`
	ChildServiceID  string                 `json:"-"`
	Extended        map[string]interface{} `json:"-"`
	EndOfLife       bool                   `json:"-"`
}

// ServiceMetadata contains metadata about the service classes
type ServiceMetadata struct { // nolint: golint
	DisplayName         string `json:"displayName,omitempty"`
	ImageURL            string `json:"imageUrl,omitempty"`
	LongDescription     string `json:"longDescription,omitempty"`
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
	DocumentationURL    string `json:"documentationUrl,omitempty"`
	SupportURL          string `json:"supportUrl,omitempty"`
}

// Service is an interface to be implemented by types that represent a single
// type of service with one or more plans
type Service interface {
	GetID() string
	GetName() string
	IsBindable() bool
	GetServiceManager() ServiceManager
	GetPlans() []Plan
	GetPlan(planID string) (Plan, bool)
	GetParentServiceID() string
	GetChildServiceID() string
	GetProperties() ServiceProperties
	IsEndOfLife() bool
}

type service struct {
	ServiceProperties
	serviceManager ServiceManager
	indexedPlans   map[string]Plan
	plans          []Plan
}

// PlanProperties represent the properties of a Plan that can be directly
// instantiated and passed to the NewPlan() constructor function which will
// carry out all necessary initialization.
type PlanProperties struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Free        bool                   `json:"free"`
	Metadata    ServicePlanMetadata    `json:"metadata,omitempty"` // nolint: lll
	Extended    map[string]interface{} `json:"-"`
	EndOfLife   bool                   `json:"-"`
	Schemas     PlanSchemas            `json:"schemas,omitempty"`
	Stability   Stability              `json:"-"`
}

// ServicePlanMetadata contains metadata about the service plans
type ServicePlanMetadata struct { // nolint: golint
	DisplayName string   `json:"displayName,omitempty"`
	Bullets     []string `json:"bullets,omitempty"`
}

// Plan is an interface to be implemented by types that represent a single
// variant or "sku" of a service
type Plan interface {
	GetID() string
	GetName() string
	GetProperties() PlanProperties
	IsEndOfLife() bool
	GetSchemas() PlanSchemas
	GetStability() Stability
}

type plan struct {
	PlanProperties
}

// NewCatalog initializes and returns a new Catalog
func NewCatalog(services []Service) Catalog {
	c := &catalog{
		services:        services,
		indexedServices: make(map[string]Service),
	}
	for _, service := range services {
		c.indexedServices[service.GetID()] = service
	}
	return c
}

// ToJSON returns a []byte containing a JSON representation of the catalog
func (c *catalog) MarshalJSON() ([]byte, error) {
	// filter out the EOL services. When MarshalJson is called
	// they won't be present in the slice
	nonEOLServices := []Service{}
	for _, service := range c.services {
		if !service.IsEndOfLife() {
			nonEOLServices = append(nonEOLServices, service)
		}
	}
	return json.Marshal(struct {
		Services []Service `json:"services"`
	}{
		Services: nonEOLServices,
	})
}

// GetServices returns all of the catalog's services
func (c *catalog) GetServices() []Service {
	return c.services
}

// GetService finds a service by serviceID in a catalog
// TODO: Test this
func (c *catalog) GetService(serviceID string) (Service, bool) {
	service, ok := c.indexedServices[serviceID]
	return service, ok
}

// NewService initialized and returns a new Service
func NewService(
	serviceProperties ServiceProperties,
	serviceManager ServiceManager,
	plans ...Plan,
) Service {
	s := service{
		ServiceProperties: serviceProperties,
		serviceManager:    serviceManager,
		plans:             plans,
		indexedPlans:      make(map[string]Plan),
	}
	for _, planIfc := range s.plans {
		p := planIfc.(plan)
		s.indexedPlans[p.GetID()] = p
	}
	return s
}

func (s service) MarshalJSON() ([]byte, error) {
	// filter out the EOL plans. When MarshalJson is called
	// they won't be present in the slice
	nonEOLPlans := []Plan{}
	for _, plan := range s.plans {
		if !plan.IsEndOfLife() {
			nonEOLPlans = append(nonEOLPlans, plan)
		}
	}
	return json.Marshal(struct {
		ServiceProperties
		Plans []Plan `json:"plans"`
	}{
		ServiceProperties: s.GetProperties(),
		Plans:             nonEOLPlans,
	})
}

func (s service) GetID() string {
	return s.ID
}

func (s service) GetName() string {
	return s.Name
}

// IsBindable returns true if a service is bindable
func (s service) IsBindable() bool {
	return s.Bindable
}

func (s service) GetServiceManager() ServiceManager {
	return s.serviceManager
}

// GetPlans returns all of the service's plans
func (s service) GetPlans() []Plan {
	return s.plans
}

// GetPlan finds a plan by planID in a service
// TODO: Test this
func (s service) GetPlan(planID string) (Plan, bool) {
	plan, ok := s.indexedPlans[planID]
	return plan, ok
}

func (s service) GetParentServiceID() string {
	return s.ParentServiceID
}

func (s service) GetChildServiceID() string {
	return s.ChildServiceID
}

func (s service) GetProperties() ServiceProperties {
	return s.ServiceProperties
}

func (s service) IsEndOfLife() bool {
	return s.EndOfLife
}

// NewPlan initializes and returns a new Plan
func NewPlan(planProperties PlanProperties) Plan {
	return plan{
		PlanProperties: planProperties,
	}
}

func (p plan) GetID() string {
	return p.ID
}

func (p plan) GetName() string {
	return p.Name
}

func (p plan) GetProperties() PlanProperties {
	return p.PlanProperties
}

func (p plan) IsEndOfLife() bool {
	return p.EndOfLife
}

func (p plan) GetSchemas() PlanSchemas {
	return p.Schemas
}

func (p plan) GetStability() Stability {
	return p.Stability
}
