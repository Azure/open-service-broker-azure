package service

import (
	"encoding/json"
	"sync"
)

// Catalog is an interface to be implemented by types that represents the
// service/plans offered by a service module or by the entire broker.
type Catalog interface {
	ToJSON() ([]byte, error)
	GetServices() []Service
	GetService(serviceID string) (Service, bool)
}

type catalog struct {
	Services        []json.RawMessage `json:"services"`
	services        []Service
	indexedServices map[string]Service
	jsonMutex       sync.Mutex
}

// ServiceProperties represent the properties of a Service that can be directly
// instantiated and passed to the NewService() constructor function which will
// carry out all necessary initialization.
type ServiceProperties struct { // nolint: golint
	Name          string           `json:"name"`
	ID            string           `json:"id"`
	Description   string           `json:"description"`
	Metadata      *ServiceMetadata `json:"metadata,omitempty"`
	Tags          []string         `json:"tags"`
	Bindable      bool             `json:"bindable"`
	PlanUpdatable bool             `json:"plan_updateable"` // Misspelling is
	// deliberate to match the spec
	ParentServiceID       string                     `json:"-"`
	ChildServiceID        string                     `json:"-"`
	ProvisionParamsSchema map[string]ParameterSchema `json:"-"`
	UpdateParamsSchema    map[string]ParameterSchema `json:"-"`
	BindingParamsSchema   map[string]ParameterSchema `json:"-"`
	Extended              map[string]interface{}     `json:"-"`
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
	ToJSON() ([]byte, error)
	GetID() string
	GetName() string
	IsBindable() bool
	GetServiceManager() ServiceManager
	GetPlans() []Plan
	GetPlan(planID string) (Plan, bool)
	GetParentServiceID() string
	GetChildServiceID() string
	GetProperties() *ServiceProperties
}

type service struct {
	*ServiceProperties
	serviceManager ServiceManager
	indexedPlans   map[string]Plan
	Plans          []json.RawMessage `json:"plans"`
	plans          []Plan
	jsonMutex      sync.Mutex
}

// PlanProperties represent the properties of a Plan that can be directly
// instantiated and passed to the NewPlan() constructor function which will
// carry out all necessary initialization.
type PlanProperties struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Free        bool                   `json:"free"`
	Metadata    *ServicePlanMetadata   `json:"metadata,omitempty"` // nolint: lll
	Extended    map[string]interface{} `json:"-"`
}

// ServicePlanMetadata contains metadata about the service plans
type ServicePlanMetadata struct { // nolint: golint
	DisplayName string   `json:"displayName,omitempty"`
	Bullets     []string `json:"bullets,omitempty"`
}

// Plan is an interface to be implemented by types that represent a single
// variant or "sku" of a service
type Plan interface {
	ToJSON() ([]byte, error)
	GetID() string
	GetName() string
	GetProperties() *PlanProperties
}

type plan struct {
	*PlanProperties
	ParameterSchemas *planSchemas `json:"schemas,omitempty"`
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
func (c *catalog) ToJSON() ([]byte, error) {
	c.jsonMutex.Lock()
	defer c.jsonMutex.Unlock()
	defer func() {
		c.Services = nil
	}()
	c.Services = []json.RawMessage{}
	for _, svc := range c.services {
		svcJSON, err := svc.ToJSON()
		if err != nil {
			return nil, err
		}
		c.Services = append(c.Services, json.RawMessage(svcJSON))
	}
	return json.Marshal(c)
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
	serviceProperties *ServiceProperties,
	serviceManager ServiceManager,
	plans ...Plan,
) Service {
	s := &service{
		ServiceProperties: serviceProperties,
		serviceManager:    serviceManager,
		plans:             plans,
		indexedPlans:      make(map[string]Plan),
	}
	paramSchemas := &planSchemas{}
	if serviceProperties.ProvisionParamsSchema != nil ||
		serviceProperties.UpdateParamsSchema != nil ||
		serviceProperties.BindingParamsSchema != nil {

		paramSchemas.addParameterSchemas(
			serviceProperties.ProvisionParamsSchema,
			serviceProperties.UpdateParamsSchema,
			serviceProperties.BindingParamsSchema,
		)
	}
	if serviceProperties.ParentServiceID == "" {
		paramSchemas.addParameterSchemas(
			getCommonProvisionParameters(),
			nil,
			nil,
		)
		if serviceProperties.ChildServiceID != "" {
			paramSchemas.addParameterSchemas(
				getParentServiceParameters(),
				nil,
				nil,
			)
		}
	} else {
		paramSchemas.addParameterSchemas(
			getChildServiceParameters(),
			nil,
			nil,
		)
	}
	for _, planIfc := range s.plans {
		p := planIfc.(*plan)
		p.ParameterSchemas = paramSchemas
		s.indexedPlans[p.GetID()] = p
	}
	return s
}

// NewServiceFromJSON returns a new Service unmarshalled from the provided
// JSON []byte
func NewServiceFromJSON(jsonBytes []byte) (Service, error) {
	s := &service{
		plans:        []Plan{},
		indexedPlans: make(map[string]Plan),
	}
	if err := json.Unmarshal(jsonBytes, s); err != nil {
		return nil, err
	}
	for _, planRawJSON := range s.Plans {
		plan, err := NewPlanFromJSON(planRawJSON)
		if err != nil {
			return nil, err
		}
		s.plans = append(s.plans, plan)
		s.indexedPlans[plan.GetID()] = plan
	}
	s.Plans = nil
	return s, nil
}

func (s *service) ToJSON() ([]byte, error) {
	s.jsonMutex.Lock()
	defer s.jsonMutex.Unlock()
	defer func() {
		s.Plans = nil
	}()
	s.Plans = []json.RawMessage{}
	for _, plan := range s.plans {
		planJSON, err := plan.ToJSON()
		if err != nil {
			return nil, err
		}
		s.Plans = append(s.Plans, json.RawMessage(planJSON))
	}
	return json.Marshal(s)
}

func (s *service) GetID() string {
	return s.ID
}

func (s *service) GetName() string {
	return s.Name
}

// IsBindable returns true if a service is bindable
func (s *service) IsBindable() bool {
	return s.Bindable
}

func (s *service) GetServiceManager() ServiceManager {
	return s.serviceManager
}

// GetPlans returns all of the service's plans
func (s *service) GetPlans() []Plan {
	return s.plans
}

// GetPlan finds a plan by planID in a service
// TODO: Test this
func (s *service) GetPlan(planID string) (Plan, bool) {
	plan, ok := s.indexedPlans[planID]
	return plan, ok
}

func (s *service) GetParentServiceID() string {
	return s.ParentServiceID
}

func (s *service) GetChildServiceID() string {
	return s.ChildServiceID
}

func (s *service) GetProperties() *ServiceProperties {
	return s.ServiceProperties
}

// NewPlan initializes and returns a new Plan
func NewPlan(planProperties *PlanProperties) Plan {
	return &plan{
		PlanProperties: planProperties,
	}
}

// NewPlanFromJSON returns a new Plan unmarshalled from the provided JSON []byte
func NewPlanFromJSON(jsonBytes []byte) (Plan, error) {
	p := &plan{}
	if err := json.Unmarshal(jsonBytes, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *plan) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *plan) GetID() string {
	return p.ID
}

func (p *plan) GetName() string {
	return p.Name
}

func (p *plan) GetProperties() *PlanProperties {
	return p.PlanProperties
}
