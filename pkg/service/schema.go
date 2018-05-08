package service

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// PlanSchemas is the root of a tree that encapsulates all plan-related schemas
// for validating input parameters to all service instance and service binding
// operations.
type PlanSchemas struct {
	ServiceInstances *InstanceSchemas `json:"service_instance,omitempty"`
	ServiceBindings  *BindingSchemas  `json:"service_binding,omitempty"`
}

// InstanceSchemas encapsulates all plan-related schemas for validating input
// paramters to all service instance operations.
type InstanceSchemas struct {
	ProvisioningParametersSchema *InputParametersSchema `json:"create,omitempty"`
	UpdatingParametersSchema     *InputParametersSchema `json:"update,omitempty"`
}

// BindingSchemas encapsulates all plan-related schemas for validating input
// parameters to all service binding operations.
type BindingSchemas struct {
	BindingParametersSchema *InputParametersSchema `json:"create,omitempty"`
}

// InputParametersSchema encapsulates schema for validating input paramaters
// to any single operation.
type InputParametersSchema struct {
	Schema             string                     `json:"$schema"`
	Type               string                     `json:"type"`
	Required           bool                       `json:"-"`
	RequiredProperties []string                   `json:"required,omitempty"`
	Properties         map[string]ParameterSchema `json:"properties,omitempty"`
	Additional         ParameterSchema            `json:"additionalProperties,omitempty"` // nolint: lll
}

// MarshalJSON defines custom JSON marshaling for InputParametersSchema and
// introduces an intermediate "parameters" property which is required by the
// OSB spec.
func (i InputParametersSchema) MarshalJSON() ([]byte, error) {
	type inputParametersSchema InputParametersSchema
	return json.Marshal(
		struct {
			Parameters inputParametersSchema `json:"parameters"`
		}{
			Parameters: inputParametersSchema(i),
		},
	)
}

// SimpleParameterSchema represents the attributes of a simple schema type
// such as a string or integer
type SimpleParameterSchema struct {
	Type          string      `json:"type"`
	Description   string      `json:"description,omitempty"`
	Required      bool        `json:"-"`
	Default       interface{} `json:"default,omitempty"`
	AllowedValues interface{} `json:"enum,omitempty"`
}

// ObjectParameterSchema represents the attributes of a complicated schema type
// that can have nested properties
type ObjectParameterSchema struct {
	Description        string                     `json:"description,omitempty"`
	Required           bool                       `json:"-"`
	RequiredProperties []string                   `json:"required,omitempty"`
	Properties         map[string]ParameterSchema `json:"properties,omitempty"`
	Additional         ParameterSchema            `json:"additionalProperties,omitempty"` // nolint: lll
}

// MarshalJSON provides functionality to marshal an
// ObjectParameterSchema to JSON
func (o ObjectParameterSchema) MarshalJSON() ([]byte, error) {
	type objectParameterSchema ObjectParameterSchema
	return json.Marshal(struct {
		Type string `json:"type"`
		objectParameterSchema
	}{
		Type: "object",
		objectParameterSchema: objectParameterSchema(o),
	})
}

// ArrayParameterSchema represents the attributes of an array type
type ArrayParameterSchema struct {
	Description string          `json:"description,omitempty"`
	Required    bool            `json:"-"`
	ItemsSchema ParameterSchema `json:"items,omitempty"`
}

// MarshalJSON provides functionality to marshal an
// ArrayParameterSchema to JSON
func (a ArrayParameterSchema) MarshalJSON() ([]byte, error) {
	type arrayParameterSchema ArrayParameterSchema
	return json.Marshal(struct {
		Type string `json:"type"`
		arrayParameterSchema
	}{
		Type:                 "array",
		arrayParameterSchema: arrayParameterSchema(a),
	})
}

// NumericParameterSchema represents a numeric type, either
// integers or floating point numbers, that can have an upper or lower bound.
type NumericParameterSchema struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Required         bool        `json:"-"`
	Default          interface{} `json:"default,omitempty"`
	Minimum          interface{} `json:"minimum,omitempty"`
	ExclusiveMinimum bool        `json:"exclusiveMinimum,omitempty"`
	Maximum          interface{} `json:"maximum,omitempty"`
	ExclusiveMaximum bool        `json:"exclusiveMaximum,omitempty"`
}

// ParameterSchema defines an interface representing a given Parameter schema.
// This could be a provision or binding parameter.
type ParameterSchema interface {
	isRequired() bool
	setRequiredProperties()
}

func (s *SimpleParameterSchema) isRequired() bool {
	return s.Required
}

func (o *ObjectParameterSchema) isRequired() bool {
	return o.Required
}

func (a *ArrayParameterSchema) isRequired() bool {
	return a.Required
}

func (n *NumericParameterSchema) isRequired() bool {
	return n.Required
}

func (o *ObjectParameterSchema) setRequiredProperties() {
	for key, param := range o.Properties {
		if param.isRequired() {
			o.RequiredProperties = append(o.RequiredProperties, key)
		}
	}
}

func (a *ArrayParameterSchema) setRequiredProperties() {
	a.ItemsSchema.setRequiredProperties()
}

// NOOP method to implement the interface
func (s *SimpleParameterSchema) setRequiredProperties() {
}

// NOOP method to implement the interface
func (n *NumericParameterSchema) setRequiredProperties() {
}

func (i *InputParametersSchema) addProperties(
	newProperties map[string]ParameterSchema,
) error {
	if newProperties == nil {
		return nil
	}
	if i.Properties == nil {
		i.Properties = make(map[string]ParameterSchema)
	}
	for key, param := range newProperties {
		param.setRequiredProperties()
		if param.isRequired() {
			i.RequiredProperties =
				append(i.RequiredProperties, key)
		}
		i.Properties[key] = param
	}
	return nil
}

func (ps *PlanSchemas) addParameterSchemas(
	instanceCreateParameters map[string]ParameterSchema,
	instanceUpdateParameters map[string]ParameterSchema,
	bindingCreateParameters map[string]ParameterSchema,
) {
	if instanceCreateParameters != nil {
		sips := ps.ServiceInstances
		if sips == nil {
			sips = &InstanceSchemas{}
			ps.ServiceInstances = sips
		}
		pps := sips.ProvisioningParametersSchema
		if pps == nil {
			pps = createEmptyParameterSchema()
			sips.ProvisioningParametersSchema = pps
		}
		err := pps.addProperties(instanceCreateParameters)
		if err != nil {
			log.Errorf("error building instance creation param schema %s", err)
		}
	}

	if instanceUpdateParameters != nil {
		sips := ps.ServiceInstances
		if sips == nil {
			sips = &InstanceSchemas{}
			ps.ServiceInstances = sips
		}
		ups := sips.UpdatingParametersSchema
		if ups == nil {
			ups = createEmptyParameterSchema()
			sips.UpdatingParametersSchema = ups
		}
		err := ups.addProperties(instanceUpdateParameters)
		log.Errorf("error building instance update param schema %s", err)
	}

	if bindingCreateParameters != nil {
		sbps := ps.ServiceBindings
		if sbps == nil {
			sbps = &BindingSchemas{}
			ps.ServiceBindings = sbps
		}
		bcps := sbps.BindingParametersSchema
		if bcps == nil {
			bcps = createEmptyParameterSchema()
			sbps.BindingParametersSchema = bcps
		}
		err := bcps.addProperties(bindingCreateParameters)
		log.Errorf("error building binding creation param schema %s", err)
	}
}

func createEmptyParameterSchema() *InputParametersSchema {
	return &InputParametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
	}
}

func getCommonProvisionParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	p["location"] = &SimpleParameterSchema{
		Type: "string",
		Description: "The Azure region in which to provision" +
			" applicable resources.",
	}
	p["resourceGroup"] = &SimpleParameterSchema{
		Type: "string",
		Description: "The (new or existing) resource group with which" +
			" to associate new resources.",
	}
	p["tags"] = &ObjectParameterSchema{
		Description: "Tags to be applied to new resources," +
			" specified as key/value pairs.",
		Additional: &SimpleParameterSchema{
			Type: "string",
		},
	}

	return p
}

func getChildServiceParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	p["parentAlias"] = &SimpleParameterSchema{
		Type: "string",
		Description: "Specifies the alias of the DBMS upon which the database " +
			"should be provisioned.",
		Required: true,
	}
	return p
}

func getParentServiceParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	p["alias"] = &SimpleParameterSchema{
		Type:        "string",
		Description: "Alias to use when provisioning databases on this DBMS",
		Required:    true,
	}
	return p
}

func (ps *PlanSchemas) addCommonSchema(sp *ServiceProperties) {
	if sp.ParentServiceID == "" {
		ps.addParameterSchemas(
			getCommonProvisionParameters(),
			nil,
			nil,
		)
		if sp.ChildServiceID != "" {
			ps.addParameterSchemas(
				getParentServiceParameters(),
				nil,
				nil,
			)
		}
	} else {
		ps.addParameterSchemas(
			getChildServiceParameters(),
			nil,
			nil,
		)
	}
}
