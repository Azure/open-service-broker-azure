package service

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
)

const jsonSchemaVersion = "http://json-schema.org/draft-04/schema#"

// PlanSchemas is the root of a tree that encapsulates all plan-related schemas
// for validating input parameters to all service instance and service binding
// operations.
type PlanSchemas struct {
	ServiceInstances InstanceSchemas `json:"service_instance,omitempty"`
	ServiceBindings  *BindingSchemas `json:"service_binding,omitempty"`
}

// InstanceSchemas encapsulates all plan-related schemas for validating input
// parameters to all service instance operations.
type InstanceSchemas struct {
	ProvisioningParametersSchema InputParametersSchema  `json:"create,omitempty"`
	UpdatingParametersSchema     *InputParametersSchema `json:"update,omitempty"`
}

// BindingSchemas encapsulates all plan-related schemas for validating input
// parameters to all service binding operations.
type BindingSchemas struct {
	BindingParametersSchema *InputParametersSchema `json:"create,omitempty"`
}

// InputParametersSchema encapsulates schema for validating input parameters
// to any single operation.
type InputParametersSchema struct {
	RequiredProperties []string                  `json:"required,omitempty"`
	PropertySchemas    map[string]PropertySchema `json:"properties,omitempty"`
	Additional         PropertySchema            `json:"additionalProperties,omitempty"` // nolint: lll
}

// MarshalJSON defines custom JSON marshaling for InputParametersSchema and
// introduces an intermediate "parameters" property which is required by the
// OSB spec.
func (i InputParametersSchema) MarshalJSON() ([]byte, error) {
	type inputParametersSchema InputParametersSchema
	type inputParametersSchemaWrapper struct {
		Schema string `json:"$schema"`
		Type   string `json:"type"`
		inputParametersSchema
	}
	return json.Marshal(
		struct {
			Parameters inputParametersSchemaWrapper `json:"parameters"`
		}{
			Parameters: inputParametersSchemaWrapper{
				Schema: jsonSchemaVersion,
				Type:   "object",
				inputParametersSchema: inputParametersSchema(i),
			},
		},
	)
}

// Validate validates the given map[string]interface{} again this schema
func (i InputParametersSchema) Validate(valMap map[string]interface{}) error {
	for _, requiredProperty := range i.RequiredProperties {
		_, ok := valMap[requiredProperty]
		if !ok {
			return NewValidationError(requiredProperty, "field is required")
		}
	}
	for k, v := range valMap {
		propertySchema, ok := i.PropertySchemas[k]
		if ok {
			if err := propertySchema.validate(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// PropertySchema is an interface for the schema of any kind of property.
type PropertySchema interface {
	validate(context string, value interface{}) error
}

// CustomPropertyValidator is a function type that describes the signature
// for functions that provide custom validation logic for property schemas.
type CustomPropertyValidator func(value interface{}) error

// StringPropertySchema represents schema for a single string property
type StringPropertySchema struct {
	Description             string                  `json:"description,omitempty"`
	MinLength               *int                    `json:"minLength,omitempty"`
	MaxLength               *int                    `json:"maxLength,omitempty"`
	AllowedValues           []string                `json:"enum,omitempty"`
	AllowedPattern          *regexp.Regexp          `json:"pattern,omitempty"`
	CustomPropertyValidator CustomPropertyValidator `json:"-"`
	DefaultValue            string                  `json:"default,omitempty"`
}

// MarshalJSON provides functionality to marshal a StringPropertySchema to JSON
func (s StringPropertySchema) MarshalJSON() ([]byte, error) {
	type stringPropertySchema StringPropertySchema
	return json.Marshal(
		struct {
			Type string `json:"type"`
			stringPropertySchema
		}{
			Type:                 "string",
			stringPropertySchema: stringPropertySchema(s),
		},
	)
}

func (s StringPropertySchema) validate(
	context string,
	value interface{},
) error {
	val, ok := value.(string)
	if !ok {
		return NewValidationError(context, "field value is not of type string")
	}
	if s.MinLength != nil && len(val) < *s.MinLength {
		return NewValidationError(
			context,
			fmt.Sprintf("field length is less than minimum %d", *s.MinLength),
		)
	}
	if s.MaxLength != nil && len(val) > *s.MaxLength {
		return NewValidationError(
			context,
			fmt.Sprintf("field length is greater than maximum %d", *s.MaxLength),
		)
	}
	if len(s.AllowedValues) > 0 {
		var found bool
		for _, allowedValue := range s.AllowedValues {
			if val == allowedValue {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(context, "field value is invalid")
		}
	}
	if s.AllowedPattern != nil {
		if !s.AllowedPattern.MatchString(val) {
			return NewValidationError(context, "field value is invalid")
		}
	}
	return nil
}

// IntPropertySchema represents schema for a single integer property
type IntPropertySchema struct {
	Description             string                  `json:"description,omitempty"`
	MinValue                *int64                  `json:"minimum,omitempty"`
	MaxValue                *int64                  `json:"maximum,omitempty"`
	AllowedValues           []int64                 `json:"enum,omitempty"`
	AllowedIncrement        *int64                  `json:"multipleOf,omitempty"`
	CustomPropertyValidator CustomPropertyValidator `json:"-"`
	DefaultValue            *int64                  `json:"default,omitempty"`
}

// MarshalJSON provides functionality to marshal an IntPropertySchema to JSON
func (i IntPropertySchema) MarshalJSON() ([]byte, error) {
	type intPropertySchema IntPropertySchema
	return json.Marshal(
		struct {
			Type string `json:"type"`
			intPropertySchema
		}{
			Type:              "integer",
			intPropertySchema: intPropertySchema(i),
		},
	)
}

func (i IntPropertySchema) validate(context string, value interface{}) error {
	floatVal, ok := value.(float64)
	if !ok {
		return NewValidationError(context, "field value is not of type int64")
	}
	val := int64(floatVal)
	if floatVal != float64(val) {
		return NewValidationError(context, "field value is not of type int64")
	}
	if i.MinValue != nil && val < *i.MinValue {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is less than minimum %d", *i.MinValue),
		)
	}
	if i.MaxValue != nil && val > *i.MaxValue {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is greater than maximum %d", *i.MaxValue),
		)
	}
	if len(i.AllowedValues) > 0 {
		var found bool
		for _, allowedValue := range i.AllowedValues {
			if val == allowedValue {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(context, "field value is invalid")
		}
	}
	if i.AllowedIncrement != nil && val%*i.AllowedIncrement != 0 {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is not a multiple of %d", *i.AllowedIncrement),
		)
	}
	return nil
}

// FloatPropertySchema represents schema for a single floating point property
type FloatPropertySchema struct {
	Description             string                  `json:"description,omitempty"`
	MinValue                *float64                `json:"minimum,omitempty"`
	MaxValue                *float64                `json:"maximum,omitempty"`
	AllowedValues           []float64               `json:"enum,omitempty"`
	AllowedIncrement        *float64                `json:"multipleOf,omitempty"`
	CustomPropertyValidator CustomPropertyValidator `json:"-"`
	DefaultValue            *float64                `json:"default,omitempty"`
}

// MarshalJSON provides functionality to marshal a FloatPropertySchema to JSON
func (f FloatPropertySchema) MarshalJSON() ([]byte, error) {
	type floatPropertySchema FloatPropertySchema
	return json.Marshal(
		struct {
			Type string `json:"type"`
			floatPropertySchema
		}{
			Type:                "number",
			floatPropertySchema: floatPropertySchema(f),
		},
	)
}

func (f FloatPropertySchema) validate(context string, value interface{}) error {
	val, ok := value.(float64)
	if !ok {
		return NewValidationError(context, "field value is not of type float64")
	}
	if f.MinValue != nil && val < *f.MinValue {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is less than minimum %f", *f.MinValue),
		)
	}
	if f.MaxValue != nil && val > *f.MaxValue {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is greater than maximum %f", *f.MaxValue),
		)
	}
	if len(f.AllowedValues) > 0 {
		var found bool
		for _, allowedValue := range f.AllowedValues {
			if val == allowedValue {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError(context, "field value is invalid")
		}
	}
	if f.AllowedIncrement != nil && math.Mod(val, *f.AllowedIncrement) != 0 {
		return NewValidationError(
			context,
			fmt.Sprintf("field value is not a multiple of %f", *f.AllowedIncrement),
		)
	}
	return nil
}

// ObjectPropertySchema represents the attributes of a complicated schema type
// that can have nested properties
type ObjectPropertySchema struct {
	Description        string                    `json:"description,omitempty"`
	RequiredProperties []string                  `json:"required,omitempty"`
	PropertySchemas    map[string]PropertySchema `json:"properties,omitempty"`
	Additional         PropertySchema            `json:"additionalProperties,omitempty"` // nolint: lll
	DefaultValue       map[string]interface{}    `json:"-"`
}

// MarshalJSON provides functionality to marshal an ObjectPropertySchema to JSON
func (o ObjectPropertySchema) MarshalJSON() ([]byte, error) {
	type objectPropertySchema ObjectPropertySchema
	return json.Marshal(struct {
		Type string `json:"type"`
		objectPropertySchema
	}{
		Type:                 "object",
		objectPropertySchema: objectPropertySchema(o),
	})
}

func (o ObjectPropertySchema) validate(
	context string,
	value interface{},
) error {
	valMap, ok := value.(map[string]interface{})
	if !ok {
		return NewValidationError(context, "field value is not of type object")
	}
	for _, requiredProperty := range o.RequiredProperties {
		_, ok := valMap[requiredProperty]
		if !ok {
			propetyContext := fmt.Sprintf("%s.%s", context, requiredProperty)
			return NewValidationError(propetyContext, "field is required")
		}
	}
	for k, v := range valMap {
		propertySchema, ok := o.PropertySchemas[k]
		if ok {
			propetyContext := fmt.Sprintf("%s.%s", context, k)
			if err := propertySchema.validate(propetyContext, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// ArrayPropertySchema represents the attributes of an array type
type ArrayPropertySchema struct {
	Description  string         `json:"description,omitempty"`
	MinItems     *int           `json:"minItems,omitempty"`
	MaxItems     *int           `json:"maxItems,omitempty"`
	ItemsSchema  PropertySchema `json:"items,omitempty"`
	DefaultValue []interface{}  `json:"-"`
}

// MarshalJSON provides functionality to marshal an
// ArrayPropertySchema to JSON
func (a ArrayPropertySchema) MarshalJSON() ([]byte, error) {
	type arrayPropertySchema ArrayPropertySchema
	return json.Marshal(struct {
		Type string `json:"type"`
		arrayPropertySchema
	}{
		Type:                "array",
		arrayPropertySchema: arrayPropertySchema(a),
	})
}

func (a ArrayPropertySchema) validate(context string, value interface{}) error {
	valArray, ok := value.([]interface{})
	if !ok {
		return NewValidationError(context, "field value is not of type array")
	}
	if a.MinItems != nil && len(valArray) < *a.MinItems {
		return NewValidationError(
			context,
			fmt.Sprintf("field contains fewer than minimum elements %d", *a.MinItems),
		)
	}
	if a.MaxItems != nil && len(valArray) > *a.MaxItems {
		return NewValidationError(
			context,
			fmt.Sprintf(
				"field contains greater than maximum elements %d",
				*a.MaxItems,
			),
		)
	}
	if a.ItemsSchema != nil {
		for i, val := range valArray {
			itemContext := fmt.Sprintf("%s[%d]", context, i)
			if err := a.ItemsSchema.validate(itemContext, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PlanSchemas) addCommonSchema(sp *ServiceProperties) {
	if p.ServiceInstances.ProvisioningParametersSchema.PropertySchemas == nil {
		p.ServiceInstances.ProvisioningParametersSchema.PropertySchemas =
			map[string]PropertySchema{}
	}
	ps := p.ServiceInstances.ProvisioningParametersSchema.PropertySchemas
	if sp.ParentServiceID == "" {
		ps["location"] = &StringPropertySchema{
			Description: "The Azure region in which to provision" +
				" applicable resources.",
		}
		ps["resourceGroup"] = &StringPropertySchema{
			Description: "The (new or existing) resource group with which" +
				" to associate new resources.",
		}
		ps["tags"] = &ObjectPropertySchema{
			Description: "Tags to be applied to new resources," +
				" specified as key/value pairs.",
			Additional: &StringPropertySchema{},
		}
		if sp.ChildServiceID != "" {
			p.ServiceInstances.ProvisioningParametersSchema.RequiredProperties =
				append(
					p.ServiceInstances.ProvisioningParametersSchema.RequiredProperties,
					"alias",
				)
			ps["alias"] = &StringPropertySchema{
				Description: "Alias to use when provisioning databases on this DBMS",
			}
		}
	} else {
		p.ServiceInstances.ProvisioningParametersSchema.RequiredProperties =
			append(
				p.ServiceInstances.ProvisioningParametersSchema.RequiredProperties,
				"parentAlias",
			)
		ps["parentAlias"] = &StringPropertySchema{
			Description: "Specifies the alias of the DBMS upon which the database " +
				"should be provisioned.",
		}
	}
}
