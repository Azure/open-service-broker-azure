package service

import (
	"encoding/json"
)

// Parameters ...
// TODO: krancour: Document this
type Parameters struct {
	Data   map[string]interface{}
	Schema *InputParametersSchema
}

// MarshalJSON ...
// TODO: krancour: Document this
func (p Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Data)
}

// UnmarshalJSON ...
// TODO: krancour: Document this
func (p *Parameters) UnmarshalJSON(bytes []byte) error {
	p.Data = map[string]interface{}{}
	return json.Unmarshal(bytes, &p.Data)
}

// GetString ...
// TODO: krancour: Document this
func (p *Parameters) GetString(key string) string {
	if p.Schema == nil {
		return ""
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return ""
	}
	stringSchema, ok := schema.(*StringPropertySchema)
	if !ok {
		return ""
	}
	valIface, ok := p.Data[key]
	if !ok {
		return stringSchema.DefaultValue
	}
	val, ok := valIface.(string)
	if !ok {
		return stringSchema.DefaultValue
	}
	return val
}

// GetInt64 ...
// TODO: krancour: Document this
func (p *Parameters) GetInt64(key string) int64 {
	if p.Schema == nil {
		return 0
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return 0
	}
	intSchema, ok := schema.(*IntPropertySchema)
	if !ok {
		return 0
	}
	valIface, ok := p.Data[key]
	if !ok {
		return *intSchema.DefaultValue
	}
	if val, ok := valIface.(*int64); ok {
		return *val
	}
	if val, ok := valIface.(int64); ok {
		return val
	}
	if val, ok := valIface.(*int32); ok {
		return int64(*val)
	}
	if val, ok := valIface.(int32); ok {
		return int64(val)
	}
	if val, ok := valIface.(*int); ok {
		return int64(*val)
	}
	if val, ok := valIface.(int); ok {
		return int64(val)
	}
	return *intSchema.DefaultValue
}

// GetFloat64 ...
// TODO: krancour: Document this
func (p *Parameters) GetFloat64(key string) float64 {
	if p.Schema == nil {
		return 0
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return 0
	}
	floatSchema, ok := schema.(*FloatPropertySchema)
	if !ok {
		return 0
	}
	valIface, ok := p.Data[key]
	if !ok {
		return *floatSchema.DefaultValue
	}
	if val, ok := valIface.(*float64); ok {
		return *val
	}
	if val, ok := valIface.(float64); ok {
		return val
	}
	if val, ok := valIface.(*float32); ok {
		return float64(*val)
	}
	if val, ok := valIface.(float32); ok {
		return float64(val)
	}
	if val, ok := valIface.(*int64); ok {
		return float64(*val)
	}
	if val, ok := valIface.(int64); ok {
		return float64(val)
	}
	if val, ok := valIface.(*int32); ok {
		return float64(*val)
	}
	if val, ok := valIface.(int32); ok {
		return float64(val)
	}
	if val, ok := valIface.(*int); ok {
		return float64(*val)
	}
	if val, ok := valIface.(int); ok {
		return float64(val)
	}
	return *floatSchema.DefaultValue
}

// GetObject ...
// TODO: krancour: Document this
func (p *Parameters) GetObject(key string) map[string]interface{} {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return nil
	}
	objectSchema, ok := schema.(*ObjectPropertySchema)
	if !ok {
		return nil
	}
	valIface, ok := p.Data[key]
	if !ok {
		return objectSchema.DefaultValue
	}
	val, ok := valIface.(map[string]interface{})
	if !ok {
		return objectSchema.DefaultValue
	}
	return val
}

// GetArray ...
// TODO: krancour: Document this
func (p *Parameters) GetArray(key string) []interface{} {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return nil
	}
	arraySchema, ok := schema.(*ArrayPropertySchema)
	if !ok {
		return nil
	}
	valIface, ok := p.Data[key]
	if !ok {
		return arraySchema.DefaultValue
	}
	val, ok := valIface.([]interface{})
	if !ok {
		return arraySchema.DefaultValue
	}
	return val
}
