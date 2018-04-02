package service

import (
	"fmt"

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
	Schema           string `json:"$schema"`
	*parameterSchema `json:",inline"`
}

type parameterSchema struct {
	Type          string                      `json:"type"`
	RequiredField bool                        `json:"-"`
	Required      []string                    `json:"required,omitempty"`
	Description   string                      `json:"description,omitempty"`
	Default       interface{}                 `json:"default,omitempty"`
	Items         *parameterSchema            `json:"items,omitempty"`
	Properties    map[string]*parameterSchema `json:"properties,omitempty"`
	Additional    *parameterSchema            `json:"additionalProperties,omitempty"` // nolint: lll
	AllowedValues interface{}                 `json:"enum,omitempty"`
}

// ParameterSchema defines an interface representing a given Parameter schema.
// This could be a provision or binding parameter.
type ParameterSchema interface {
	SetRequired(bool)
	IsRequired() bool
	AddParameters(map[string]ParameterSchema) error
	SetAdditionalPropertiesType(kind string)
	SetDefault(defaultVal interface{})
	SetItems(items ParameterSchema) error
	SetAllowedValues(interface{})
}

func (p *parameterSchema) SetRequired(isRequired bool) {
	p.RequiredField = isRequired
}

func (p *parameterSchema) IsRequired() bool {
	return p.RequiredField
}

func (p *parameterSchema) AddParameters(
	newParams map[string]ParameterSchema,
) error {
	if newParams == nil {
		return nil
	}
	if p.Properties == nil {
		p.Properties = make(map[string]*parameterSchema)
	}
	for key, param := range newParams {
		ps, ok := param.(*parameterSchema)
		if !ok {
			return fmt.Errorf("Unknown parameters object")
		}
		p.Properties[key] = ps
		if param.IsRequired() {
			p.Required = append(p.Required, key)
		}
	}
	return nil
}

func (p *parameterSchema) SetAdditionalPropertiesType(kind string) {
	p.Additional = &parameterSchema{
		Type: kind,
	}
}

func (p *parameterSchema) SetDefault(defaultVal interface{}) {
	p.Default = defaultVal
}

func (p *parameterSchema) SetAllowedValues(allowedValues interface{}) {
	p.AllowedValues = allowedValues
}

func (p *parameterSchema) SetItems(itemSchema ParameterSchema) error {
	items, ok := itemSchema.(*parameterSchema)
	if !ok {
		return fmt.Errorf("Unknown parameters object")
	}
	p.Items = items
	return nil
}

// NewParameterSchema returns an instance of a type that implements
// the ParameterSchema interface. Service authors can use the interface
// methods to fully construct the object
func NewParameterSchema(
	typeString string,
	description string,
) ParameterSchema {
	ps := &parameterSchema{
		Type:        typeString,
		Description: description,
	}
	return ps
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
		err := pps.Parameters.AddParameters(instanceCreateParameters)
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
		err := ups.Parameters.AddParameters(instanceUpdateParameters)
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
		err := bcps.Parameters.AddParameters(bindingCreateParameters)
		log.Errorf("error building binding creation param schema %s", err)
	}
}

func createEmptyParameterSchema() *parametersSchema {
	p := &parametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		parameterSchema: &parameterSchema{
			Type: "object",
		},
	}
	return p
}

func getCommonProvisionParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	p["location"] = NewParameterSchema(
		"string",
		"The Azure region in which to provision applicable resources.",
	)

	p["resourceGroup"] = NewParameterSchema(
		"string",
		"The (new or existing) resource group with which to associate new"+
			" resources.",
	)

	tagsSchema := NewParameterSchema(
		"object",
		"Tags to be applied to new resources, specified as key/value pairs.",
	)
	tagsSchema.SetAdditionalPropertiesType("string")

	p["tags"] = tagsSchema

	return p
}
