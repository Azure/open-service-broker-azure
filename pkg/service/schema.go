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
	SetAdditionalPropertiesType(kind string)
	SetDefault(defaultVal interface{})
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

// NewSimpleParameterSchema returns an instance of a type that implements
// the ParameterSchema interface with a type and description.
func NewSimpleParameterSchema(
	typeString string,
	description string,
) ParameterSchema {
	ps := &parameterSchema{
		Type:        typeString,
		Description: description,
	}
	return ps
}

// NewObjectParameterSchema returns an instance of a type that implements
// the ParameterSchema interface. This assumes the "object" type and
// adds any specified properties to the schema instance.
func NewObjectParameterSchema(
	description string,
	properties map[string]ParameterSchema,
) ParameterSchema {
	ps := &parameterSchema{
		Type:        "object",
		Description: description,
	}

	err := ps.AddParameters(properties)
	if err != nil {
		log.Errorf(
			"error adding parameter properties to object schema: %s",
			err,
		)
	}
	return ps
}

// NewArrayParameterSchema returns an instance of a type that implements
// the ParameterSchema interface. This assumes the "object" type and
// adds any specified properties to the schema instance.
func NewArrayParameterSchema(
	description string,
	itemSchema ParameterSchema,
) ParameterSchema {
	ps := &parameterSchema{
		Type:        "array",
		Description: description,
	}

	err := ps.SetItems(itemSchema)
	if err != nil {
		log.Errorf("Error adding items to the array schema: %s", err)
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
	p["location"] = NewSimpleParameterSchema(
		"string",
		"The Azure region in which to provision applicable resources.",
	)

	p["resourceGroup"] = NewSimpleParameterSchema(
		"string",
		"The (new or existing) resource group with which to associate new"+
			" resources.",
	)

	tagsSchema := NewObjectParameterSchema(
		"Tags to be applied to new resources, specified as key/value pairs.",
		nil,
	)
	tagsSchema.SetAdditionalPropertiesType("string")

	p["tags"] = tagsSchema

	return p
}

func getChildServiceParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	parentAliasSchema := NewSimpleParameterSchema(
		"string",
		"Specifies the alias of the DBMS upon which the database "+
			"should be provisioned.",
	)
	parentAliasSchema.SetRequired(true)
	p["parentAlias"] = parentAliasSchema
	return p
}

func getParentServiceParameters() map[string]ParameterSchema {
	p := map[string]ParameterSchema{}
	aliasSchema := NewSimpleParameterSchema(
		"string",
		"Alias to use when provisioning databases on this DBMS",
	)
	aliasSchema.SetRequired(true)
	p["alias"] = aliasSchema
	return p
}
