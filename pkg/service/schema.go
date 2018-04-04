package service

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type planSchemas struct {
	ServiceInstances *instanceSchemas `json:"service_instance,omitempty"`
	ServiceBindings  *bindingSchemas  `json:"service_binding,omitempty"`
}

type instanceSchemas struct {
	ProvisioningParametersSchema *provisioningParametersSchema `json:"create,omitempty"` // nolint: lll
	UpdatingParametersSchema     *updatingParametersSchema     `json:"update,omitempty"` // nolint: lll
}

type provisioningParametersSchema struct {
	Parameters *parametersSchema `json:"parameters,omitempty"`
}

type updatingParametersSchema struct {
	Parameters *parametersSchema `json:"parameters,omitempty"`
}
type bindingSchemas struct {
	BindingParametersSchema *bindingSchema `json:"create,omitempty"`
}

type bindingSchema struct {
	Parameters *parametersSchema `json:"parameters,omitempty"`
}

type parametersSchema struct {
	Schema             string                     `json:"$schema"`
	Type               string                     `json:"type"`
	Description        string                     `json:"description,omitempty"`
	Required           bool                       `json:"-"`
	RequiredProperties []string                   `json:"required,omitempty"`
	Properties         map[string]ParameterSchema `json:"properties,omitempty"`
	Additional         ParameterSchema            `json:"additionalProperties,omitempty"` // nolint: lll
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
	Description        string
	Required           bool
	requiredProperties []string
	Properties         map[string]ParameterSchema
	Additional         ParameterSchema
}

// MarshalJSON provides functionality to marshal an
// ObjectParameterSchema to JSON
func (o *ObjectParameterSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type               string                     `json:"type"`
		Description        string                     `json:"description,omitempty"`
		RequiredProperties []string                   `json:"required,omitempty"`
		Properties         map[string]ParameterSchema `json:"properties,omitempty"`
		Additional         ParameterSchema            `json:"additionalProperties,omitempty"` // nolint: lll
	}{
		Type:               "object",
		Description:        o.Description,
		RequiredProperties: o.requiredProperties,
		Properties:         o.Properties,
		Additional:         o.Additional,
	})
}

// ArrayParameterSchema represents the attributes of an array type
type ArrayParameterSchema struct {
	Description string
	Required    bool
	ItemsSchema ParameterSchema
}

// MarshalJSON provides functionality to marshal an
// ArrayParameterSchema to JSON
func (a *ArrayParameterSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type        string          `json:"type"`
		Description string          `json:"description,omitempty"`
		ItemsSchema ParameterSchema `json:"items,omitempty"`
	}{
		Type:        "array",
		Description: a.Description,
		ItemsSchema: a.ItemsSchema,
	})
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

// IsRequired returns true if the property is required
func (a *ArrayParameterSchema) isRequired() bool {
	return a.Required
}

func (o *ObjectParameterSchema) setRequiredProperties() {
	for key, param := range o.Properties {
		if param.isRequired() {
			o.requiredProperties = append(o.requiredProperties, key)
		}
	}
}

func (a *ArrayParameterSchema) setRequiredProperties() {
	a.ItemsSchema.setRequiredProperties()
}

// NOOP method to implement the interface
func (s *SimpleParameterSchema) setRequiredProperties() {
}

func (ps *planSchemas) setRequiredProperties() {
	if ps.ServiceInstances != nil {
		sips := ps.ServiceInstances
		if sips.ProvisioningParametersSchema != nil {
			provisionSchema := sips.ProvisioningParametersSchema.Parameters
			for key, param := range provisionSchema.Properties {
				param.setRequiredProperties()
				if param.isRequired() {
					provisionSchema.RequiredProperties =
						append(provisionSchema.RequiredProperties, key)
				}
			}
		}
		if sips.UpdatingParametersSchema != nil {
			updatingSchema := sips.UpdatingParametersSchema.Parameters
			for key, param := range updatingSchema.Properties {
				param.setRequiredProperties()
				if param.isRequired() {
					updatingSchema.RequiredProperties =
						append(updatingSchema.RequiredProperties, key)
				}
			}
		}
	}
	if ps.ServiceBindings != nil {
		bps := ps.ServiceBindings.BindingParametersSchema
		if bps.Parameters != nil {
			for key, param := range bps.Parameters.Properties {
				param.setRequiredProperties()
				if param.isRequired() {
					bps.Parameters.RequiredProperties =
						append(bps.Parameters.RequiredProperties, key)
				}
			}
		}
	}
}

func (p *parametersSchema) addProperties(
	newProperties map[string]ParameterSchema,
) error {
	if newProperties == nil {
		return nil
	}
	if p.Properties == nil {
		p.Properties = make(map[string]ParameterSchema)
	}
	for key, param := range newProperties {
		p.Properties[key] = param
	}
	return nil
}

func (ps *planSchemas) addParameterSchemas(
	instanceCreateParameters map[string]ParameterSchema,
	instanceUpdateParameters map[string]ParameterSchema,
	bindingCreateParameters map[string]ParameterSchema,
) {
	if instanceCreateParameters != nil {
		sips := ps.ServiceInstances
		if sips == nil {
			sips = &instanceSchemas{}
			ps.ServiceInstances = sips
		}
		pps := sips.ProvisioningParametersSchema
		if pps == nil {
			pps = &provisioningParametersSchema{
				Parameters: createEmptyParameterSchema(),
			}
			sips.ProvisioningParametersSchema = pps
		}
		err := pps.Parameters.addProperties(instanceCreateParameters)
		if err != nil {
			log.Errorf("error building instance creation param schema %s", err)
		}
	}

	if instanceUpdateParameters != nil {
		sips := ps.ServiceInstances
		if sips == nil {
			sips = &instanceSchemas{}
			ps.ServiceInstances = sips
		}
		ups := sips.UpdatingParametersSchema
		if ups == nil {
			ups = &updatingParametersSchema{
				Parameters: createEmptyParameterSchema(),
			}
			sips.UpdatingParametersSchema = ups
		}
		err := ups.Parameters.addProperties(instanceUpdateParameters)
		log.Errorf("error building instance update param schema %s", err)
	}
	if bindingCreateParameters != nil {

		sbps := ps.ServiceBindings
		if sbps == nil {
			sbps = &bindingSchemas{}
			ps.ServiceBindings = sbps
		}
		bcps := sbps.BindingParametersSchema
		if bcps == nil {
			bcps = &bindingSchema{
				Parameters: createEmptyParameterSchema(),
			}
			sbps.BindingParametersSchema = bcps
		}
		err := bcps.Parameters.addProperties(bindingCreateParameters)
		log.Errorf("error building binding creation param schema %s", err)
	}
}

func createEmptyParameterSchema() *parametersSchema {
	p := &parametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
	}
	return p
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
