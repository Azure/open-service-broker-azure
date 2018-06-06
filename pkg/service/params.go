package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/slice"
)

// Parameters ...
// TODO: krancour: Document this
type Parameters struct {
	Codec  crypto.Codec
	Schema *InputParametersSchema
	Data   map[string]interface{}
}

// MarshalJSON ...
// TODO: krancour: Document this
func (p Parameters) MarshalJSON() ([]byte, error) {
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Schema == nil {
		return nil, errors.New(
			`error marshaling parameters: cannot marshal without a schema`,
		)
	}
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Codec == nil {
		return nil, errors.New(
			`error marshaling parameters: cannot marshal without a codec`,
		)
	}
	data := map[string]interface{}{}
	for k, schema := range p.Schema.PropertySchemas {
		if v, ok := p.Data[k]; ok {
			if slice.ContainsString(p.Schema.SecureProperties, k) {
				if _, ok := schema.(*StringPropertySchema); !ok {
					return nil, fmt.Errorf(
						`error marshaling parameters: cannot encrypt non-string field "%s"`,
						k,
					)
				}
				vStr, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf(
						`error marshaling parameters: cannot encrypt non-string value of `+
							`string field "%s"`,
						k,
					)
				}
				vBytes := []byte(vStr)
				vBytes, err := p.Codec.Encrypt(vBytes)
				if err != nil {
					return nil, err
				}
				v = string(vBytes)
			}
			data[k] = v
		}
	}
	return json.Marshal(data)
}

// UnmarshalJSON ...
// TODO: krancour: Document this
func (p *Parameters) UnmarshalJSON(bytes []byte) error {
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Schema == nil {
		return errors.New(
			`error marshaling parameters: cannot unmarshal without a schema`,
		)
	}
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Codec == nil {
		return errors.New(
			`error marshaling parameters: cannot unmarshal without a codec`,
		)
	}
	p.Data = map[string]interface{}{}
	data := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	for k, schema := range p.Schema.PropertySchemas {
		if v, ok := data[k]; ok {
			if slice.ContainsString(p.Schema.SecureProperties, k) {
				if _, ok := schema.(*StringPropertySchema); !ok {
					return fmt.Errorf(
						`error unmarshaling parameters: cannot decrypt non-string field `+
							`"%s"`,
						k,
					)
				}
				vStr, ok := v.(string)
				if !ok {
					return fmt.Errorf(
						`error marshaling parameters: cannot decrypt non-string value of `+
							`string field "%s"`,
						k,
					)
				}
				vBytes := []byte(vStr)
				vBytes, err := p.Codec.Decrypt(vBytes)
				if err != nil {
					return err
				}
				v = string(vBytes)
			}
			p.Data[k] = v
		}
	}
	return nil
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
		if intSchema.DefaultValue == nil {
			return 0
		}
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
	if intSchema.DefaultValue == nil {
		return 0
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
		if floatSchema.DefaultValue == nil {
			return 0
		}
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
	if floatSchema.DefaultValue == nil {
		return 0
	}
	return *floatSchema.DefaultValue
}

// GetObject ...
// TODO: krancour: Document this
func (p *Parameters) GetObject(key string) map[string]interface{} {
	if p.Schema == nil {
		return map[string]interface{}{}
	}
	schema, ok := p.Schema.PropertySchemas[key]
	if !ok {
		return map[string]interface{}{}
	}
	objectSchema, ok := schema.(*ObjectPropertySchema)
	if !ok {
		return map[string]interface{}{}
	}
	valIface, ok := p.Data[key]
	if !ok {
		if objectSchema.DefaultValue == nil {
			return map[string]interface{}{}
		}
		return objectSchema.DefaultValue
	}
	val, ok := valIface.(map[string]interface{})
	if !ok {
		if objectSchema.DefaultValue == nil {
			return map[string]interface{}{}
		}
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
