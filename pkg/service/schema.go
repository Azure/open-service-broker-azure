package service

// ParameterSchemas is a set of optional JSONSchemas that describe
// the expected parameters for creation and update of instances and
// creation of bindings.
type ParameterSchemas struct {
	ServiceInstances *InstanceSchema `json:"service_instance,omitempty"`
	ServiceBindings  *BindingSchema  `json:"service_binding,omitempty"`
}

// InstanceSchema represents a plan's schemas for creation and
// update of an API resource.
type InstanceSchema struct {
	Create *InputParameters `json:"create,omitempty"`
	Update *InputParameters `json:"update,omitempty"`
}

// BindingSchema represents a plan's schemas for the parameters
// accepted for binding creation.
type BindingSchema struct {
	Create *InputParameters `json:"create,omitempty"`
}

// InputParameters represents a schema for input parameters for creation or
// update of an API resource.
type InputParameters struct {
	Parameters *ParametersSchema `json:"parameters,omitempty"`
}

// ParametersSchema represents the individual schema for a service
type ParametersSchema struct {
	Schema     string                 `json:"$schema"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// Parameter represents the individual schema for a given parameter
type Parameter struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Items       interface{}            `json:"items,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Additional  interface{}            `json:"additionalProperties,omitempty"`
}

// GetEmptyParameterSchema builds a stub Parameters object for use in schema
// definition
func GetEmptyParameterSchema() *ParametersSchema {
	p := &ParametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
	}
	return p
}

// GetCommonSchema builds a default schema object with
// location, resource group and tags
func GetCommonSchema() *ParametersSchema {
	p := GetEmptyParameterSchema()
	props := map[string]interface{}{}
	props["location"] = Parameter{
		Type: "string",
		Description: "The Azure region in which to provision " +
			"applicable resources.",
	}
	props["resourceGroup"] = Parameter{
		Type: "string",
		Description: "The (new or existing) resource group with " +
			"which to associate new resources.",
	}
	props["tags"] = Parameter{
		Type: "object",
		Description: "Tags to be applied to new resources, specified " +
			"as key/value pairs.",
		Additional: Parameter{
			Type: "string",
		},
	}
	p.Properties = props
	return p
}
