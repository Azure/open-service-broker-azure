package service

import "encoding/json"
import "sync"

// Catalog is an interface to be implemented by types that represents the
// service/plans offered by a service module or by the entire broker.
type Catalog interface {
	ToJSONString() (string, error)
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
	ToJSONString() (string, error)
	GetID() string
	GetPlans() []Plan
	GetPlan(planID string) (Plan, bool)
}

type service struct {
	*ServiceProperties
	indexedPlans map[string]Plan
	Plans        []json.RawMessage `json:"plans"`
	plans        []Plan
	jsonMutex    sync.Mutex
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
	ToJSONString() (string, error)
	GetID() string
}

type plan struct {
	*PlanProperties
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

// NewCatalogFromJSONString returns a new Catalog unmarshalled from the
// provided JSON string
func NewCatalogFromJSONString(jsonStr string) (Catalog, error) {
	c := &catalog{
		services:        []Service{},
		indexedServices: make(map[string]Service),
	}
	err := json.Unmarshal([]byte(jsonStr), c)
	if err != nil {
		return nil, err
	}
	for _, svcRawJSON := range c.Services {
		svc, err := NewServiceFromJSONString(string(svcRawJSON))
		if err != nil {
			return nil, err
		}
		c.services = append(c.services, svc)
		c.indexedServices[svc.GetID()] = svc
	}
	c.Services = nil
	return c, nil
}

// ToJSONString returns a string containing a JSON representation of the
// catalog
func (c *catalog) ToJSONString() (string, error) {
	c.jsonMutex.Lock()
	defer c.jsonMutex.Unlock()
	defer func() {
		c.Services = nil
	}()
	c.Services = []json.RawMessage{}
	for _, svc := range c.services {
		svcJSON, err := svc.ToJSONString()
		if err != nil {
			return "", err
		}
		c.Services = append(c.Services, json.RawMessage(svcJSON))
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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
func NewService(serviceProperties *ServiceProperties, plans ...Plan) Service {
	s := &service{
		ServiceProperties: serviceProperties,
		plans:             plans,
		indexedPlans:      make(map[string]Plan),
	}
	for _, plan := range s.plans {
		s.indexedPlans[plan.GetID()] = plan
	}
	return s
}

// NewServiceFromJSONString returns a new Service unmarshalled from the provided
// JSON
func NewServiceFromJSONString(jsonStr string) (Service, error) {
	s := &service{
		plans:        []Plan{},
		indexedPlans: make(map[string]Plan),
	}
	err := json.Unmarshal([]byte(jsonStr), s)
	if err != nil {
		return nil, err
	}
	for _, planRawJSON := range s.Plans {
		plan, err := NewPlanFromJSONString(string(planRawJSON))
		if err != nil {
			return nil, err
		}
		s.plans = append(s.plans, plan)
		s.indexedPlans[plan.GetID()] = plan
	}
	s.Plans = nil
	return s, nil
}

func (s *service) ToJSONString() (string, error) {
	s.jsonMutex.Lock()
	defer s.jsonMutex.Unlock()
	defer func() {
		s.Plans = nil
	}()
	s.Plans = []json.RawMessage{}
	for _, plan := range s.plans {
		planJSON, err := plan.ToJSONString()
		if err != nil {
			return "", err
		}
		s.Plans = append(s.Plans, json.RawMessage(planJSON))
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *service) GetID() string {
	return s.ID
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

// NewPlan initializes and returns a new Plan
func NewPlan(planProperties *PlanProperties) Plan {
	return &plan{
		planProperties,
	}
}

// NewPlanFromJSONString returns a new Plan unmarshalled from the provided JSON
func NewPlanFromJSONString(jsonStr string) (Plan, error) {
	p := &plan{}
	err := json.Unmarshal([]byte(jsonStr), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *plan) ToJSONString() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *plan) GetID() string {
	return p.ID
}
