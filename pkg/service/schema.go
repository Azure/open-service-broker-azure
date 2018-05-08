package service

import (
	"encoding/json"
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
// paramters to all service instance operations.
type InstanceSchemas struct {
	ProvisioningParametersSchema InputParametersSchema  `json:"create,omitempty"`
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
	RequiredProperties []string                  `json:"required,omitempty"`
	Properties         map[string]PropertySchema `json:"properties,omitempty"`
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

// SimplePropertySchema represents the attributes of a simple schema type
// such as a string or integer
type SimplePropertySchema struct {
	Type          string      `json:"type"`
	Description   string      `json:"description,omitempty"`
	Default       interface{} `json:"default,omitempty"`
	AllowedValues interface{} `json:"enum,omitempty"`
}

// ObjectPropertySchema represents the attributes of a complicated schema type
// that can have nested properties
type ObjectPropertySchema struct {
	Description        string                    `json:"description,omitempty"`
	RequiredProperties []string                  `json:"required,omitempty"`
	Properties         map[string]PropertySchema `json:"properties,omitempty"`
	Additional         PropertySchema            `json:"additionalProperties,omitempty"` // nolint: lll
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

// ArrayPropertySchema represents the attributes of an array type
type ArrayPropertySchema struct {
	Description string         `json:"description,omitempty"`
	ItemsSchema PropertySchema `json:"items,omitempty"`
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

// NumericPropertySchema represents a numeric type, either
// integers or floating point numbers, that can have an upper or lower bound.
type NumericPropertySchema struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Default          interface{} `json:"default,omitempty"`
	Minimum          interface{} `json:"minimum,omitempty"`
	ExclusiveMinimum bool        `json:"exclusiveMinimum,omitempty"`
	Maximum          interface{} `json:"maximum,omitempty"`
	ExclusiveMaximum bool        `json:"exclusiveMaximum,omitempty"`
}

// PropertySchema is an interface for the schema of any kind of property.
type PropertySchema interface{}

func (ps *PlanSchemas) addCommonSchema(sp *ServiceProperties) {
	if ps.ServiceInstances.ProvisioningParametersSchema.Properties == nil {
		ps.ServiceInstances.ProvisioningParametersSchema.Properties = map[string]PropertySchema{}
	}
	if sp.ParentServiceID == "" {
		ps.ServiceInstances.ProvisioningParametersSchema.Properties["location"] = &SimplePropertySchema{
			Type: "string",
			Description: "The Azure region in which to provision" +
				" applicable resources.",
		}
		ps.ServiceInstances.ProvisioningParametersSchema.Properties["resourceGroup"] = &SimplePropertySchema{
			Type: "string",
			Description: "The (new or existing) resource group with which" +
				" to associate new resources.",
		}
		ps.ServiceInstances.ProvisioningParametersSchema.Properties["tags"] = &ObjectPropertySchema{
			Description: "Tags to be applied to new resources," +
				" specified as key/value pairs.",
			Additional: &SimplePropertySchema{
				Type: "string",
			},
		}
		if sp.ChildServiceID != "" {
			ps.ServiceInstances.ProvisioningParametersSchema.RequiredProperties =
				append(ps.ServiceInstances.ProvisioningParametersSchema.RequiredProperties, "alias")
			ps.ServiceInstances.ProvisioningParametersSchema.Properties["alias"] = &SimplePropertySchema{
				Type:        "string",
				Description: "Alias to use when provisioning databases on this DBMS",
			}
		}
	} else {
		ps.ServiceInstances.ProvisioningParametersSchema.RequiredProperties =
			append(ps.ServiceInstances.ProvisioningParametersSchema.RequiredProperties, "parentAlias")
		ps.ServiceInstances.ProvisioningParametersSchema.Properties["parentAlias"] = &SimplePropertySchema{
			Type: "string",
			Description: "Specifies the alias of the DBMS upon which the database " +
				"should be provisioned.",
		}
	}
}
