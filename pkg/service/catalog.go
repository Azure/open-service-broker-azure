package service

import "encoding/json"

// Catalog is an interface to be implemented by types that represents the
// service/plans offered by a service module or by the entire broker.
type Catalog interface {
	ToJSONString() (string, error)
	GetServices() []Service
	GetService(serviceID string) (Service, bool)
}

type catalog struct {
	Services        []Service `json:"services"`
	indexedServices map[string]Service
}

// ServiceProperties represent the properties of a Service that can be directly
// instantiated and passed to the NewService() constructor function which will
// carry out all necessary initialization.
type ServiceProperties struct {
	Name          string   `json:"name"`
	ID            string   `json:"id"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Bindable      bool     `json:"bindable"`
	PlanUpdatable bool     `json:"plan_updateable"` // Mispelling is deliberate to match the spec
}

// Service is an interface to be implemented by types that represent a single
// type of service with one or more plans
type Service interface {
	GetID() string
	GetPlan(planID string) (Plan, bool)
}

type service struct {
	*ServiceProperties
	indexedPlans map[string]Plan
	Plans        []Plan `json:"plans"`
}

// PlanProperties represent the properties of a Plan that can be directly
// instantiated and passed to the NewPlan() constructor function which will
// carry out all necessary initialization.
type PlanProperties struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Free        bool   `json:"free"`
}

// Plan is an interface to be implemented by types that represent a single
// variant or "sku" of a service
type Plan interface {
	GetID() string
}

type plan struct {
	*PlanProperties
}

// NewCatalog initializes and returns a new Catalog
func NewCatalog(services []Service) Catalog {
	c := &catalog{
		Services:        services,
		indexedServices: make(map[string]Service),
	}
	for _, service := range services {
		c.indexedServices[service.GetID()] = service
	}
	return c
}

// ToJSONString returns a string containing a JSON representation of the
// catalog
func (c *catalog) ToJSONString() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetServices returns all of the catalog's services
func (c *catalog) GetServices() []Service {
	return c.Services
}

// GetService finds a service by serviceID in a catalog
// TODO: Test this
func (c *catalog) GetService(serviceID string) (Service, bool) {
	service, ok := c.indexedServices[serviceID]
	return service, ok
}

// NewService initialized and returns a new Service
func NewService(serviceProperties *ServiceProperties, plans ...Plan) Service {
	s := &service{
		ServiceProperties: serviceProperties,
		Plans:             plans,
		indexedPlans:      make(map[string]Plan),
	}
	for _, plan := range s.Plans {
		s.indexedPlans[plan.GetID()] = plan
	}
	return s
}

func (s *service) GetID() string {
	return s.ID
}

// GetPlan finds a plan by planID in a service
// TODO: Test this
func (s *service) GetPlan(planID string) (Plan, bool) {
	plan, ok := s.indexedPlans[planID]
	return plan, ok
}

// NewPlan initializes and returns a new Plan
func NewPlan(planProperties *PlanProperties) Plan {
	return &plan{
		planProperties,
	}
}

func (p *plan) GetID() string {
	return p.ID
}
