package service

// ParameterSchemas is a set of optional JSONSchemas that describe
// the expected parameters for creation and update of instances and
// creation of bindings.
type ParameterSchemas struct {
	ServiceInstances *instancesSchema `json:"service_instance,omitempty"`
	ServiceBindings  *BindingsSchema  `json:"service_binding,omitempty"`
}

type instancesSchema struct {
	Create *ProvisioningParametersSchema `json:"create,omitempty"`
	Update *UpdatingParametersSchema     `json:"update,omitempty"`
}

// ProvisioningParametersSchema represents the schema for any parameters that
// might be needed for provisioning a service instance
type ProvisioningParametersSchema struct {
	Parameters *ParametersSchema `json:"parameters,omitempty"`
}

// UpdatingParametersSchema represents the schema for any parameters that
// might be needed for updating a service instance
type UpdatingParametersSchema struct {
	Parameters *ParametersSchema `json:"parameters,omitempty"`
}

// BindingsSchema represents a plan's schemas for the parameters
// accepted for binding creation.
type BindingsSchema struct {
	Create *BindingSchema `json:"create,omitempty"`
}

// BindingSchema represents the schema for any parameters that
// might be needed for creating a binding
type BindingSchema struct {
	Parameters *ParametersSchema `json:"parameters,omitempty"`
}

// ParametersSchema represents the individual schema for a service
type ParametersSchema struct {
	Schema     string               `json:"$schema"`
	Type       string               `json:"type"`
	Properties map[string]Parameter `json:"properties,omitempty"`
	Required   []string             `json:"required,omitempty"`
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
// definition. This can be used to construct provision parameter,
// update paramater or binding paramater schemas.
func GetEmptyParameterSchema() *ParametersSchema {
	p := &ParametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
	}
	return p
}

// GetCommonProvisionParametersSchema builds a default schema object with
// location, resource group and tags
func GetCommonProvisionParametersSchema() *ParametersSchema {
	p := GetEmptyParameterSchema()
	props := map[string]Parameter{}
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
