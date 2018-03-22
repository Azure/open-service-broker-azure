package service

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
	Schema     string                      `json:"$schema"`
	Type       string                      `json:"type"`
	Properties map[string]*ParameterSchema `json:"properties,omitempty"`
	Required   []string                    `json:"required,omitempty"`
}

// ParameterSchema represents the individual schema for a given parameter
type ParameterSchema struct {
	Type        string                      `json:"type"`
	Required    bool                        `json:"-"`
	Description string                      `json:"description,omitempty"`
	Default     interface{}                 `json:"default,omitempty"`
	Items       *ParameterSchema            `json:"items,omitempty"`
	Properties  map[string]*ParameterSchema `json:"properties,omitempty"`
	Additional  *ParameterSchema            `json:"additionalProperties,omitempty"`
}

func (s *parametersSchema) addParameters(
	newParams map[string]*ParameterSchema,
) {
	if s.Properties == nil {
		s.Properties = newParams
		for key, param := range newParams {
			if param.Required {
				s.Required = append(s.Required, key)
			}
		}
	} else {
		for key, param := range newParams {
			s.Properties[key] = param
			if param.Required {
				s.Required = append(s.Required, key)
			}
		}
	}
}

func (ps *planSchemas) addParameters(
	instanceCreateParameters map[string]*ParameterSchema,
	instanceUpdateParameters map[string]*ParameterSchema,
	bindingCreateParameters map[string]*ParameterSchema,
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
		pps.Parameters.addParameters(instanceCreateParameters)
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
		ups.Parameters.addParameters(instanceUpdateParameters)
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
		bcps.Parameters.addParameters(bindingCreateParameters)
	}
}

func createEmptyParameterSchema() *parametersSchema {
	p := &parametersSchema{
		Schema: "http://json-schema.org/draft-04/schema#",
		Type:   "object",
	}
	return p
}

func getCommonProvisionParameters() map[string]*ParameterSchema {
	p := map[string]*ParameterSchema{}
	p["location"] = &ParameterSchema{
		Type: "string",
		Description: "The Azure region in which to provision " +
			"applicable resources.",
	}
	p["resourceGroup"] = &ParameterSchema{
		Type: "string",
		Description: "The (new or existing) resource group with " +
			"which to associate new resources.",
	}
	p["tags"] = &ParameterSchema{
		Type: "object",
		Description: "Tags to be applied to new resources, specified " +
			"as key/value pairs.",
		Additional: &ParameterSchema{
			Type: "string",
		},
	}
	return p
}
