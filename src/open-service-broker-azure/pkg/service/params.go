package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"open-service-broker-azure/pkg/crypto"
	"open-service-broker-azure/pkg/slice"
)

// Parameters is a wrapper for a map that uses a schema to inform data access
// and both schema and a codec to effect marshaling and unmarshaling with
// seamless encryption and decryption of sensitive string fields.
type Parameters struct {
	Codec  crypto.Codec
	Schema KeyedPropertySchemaContainer
	Data   map[string]interface{}
}

// MarshalJSON marshals Parameters to JSON
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
	ips, ok := p.Schema.(*InputParametersSchema)
	if !ok {
		return nil, errors.New(
			`error marshaling parameters: cannot marshal with a schema that is ` +
				`not an *InputParametersSchema`,
		)
	}
	data := map[string]interface{}{}
	for k, schema := range ips.PropertySchemas {
		if v, ok := p.Data[k]; ok {
			if slice.ContainsString(ips.SecureProperties, k) {
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

// UnmarshalJSON umarshals JSON to a Parameters objects
func (p *Parameters) UnmarshalJSON(bytes []byte) error {
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Schema == nil {
		return errors.New(
			`error unmarshaling parameters: cannot unmarshal without a schema`,
		)
	}
	// TODO: krancour: Ideally, if we constrain how Params are created using some
	// constructor-like function, perhaps we can forgo this check.
	if p.Codec == nil {
		return errors.New(
			`error unmarshaling parameters: cannot unmarshal without a codec`,
		)
	}
	ips, ok := p.Schema.(*InputParametersSchema)
	if !ok {
		return errors.New(
			`error unmarshaling parameters: cannot unmarshal with a schema that is ` +
				`not an *InputParametersSchema`,
		)
	}
	p.Data = map[string]interface{}{}
	if ips != nil {
		data := map[string]interface{}{}
		if err := json.Unmarshal(bytes, &data); err != nil {
			return err
		}
		for k, schema := range ips.PropertySchemas {
			if v, ok := data[k]; ok {
				if slice.ContainsString(ips.SecureProperties, k) {
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
	}
	return nil
}

// GetString retrieves a string by key from the Parameters' underlying map. If
// the key does not exist in the schema, an empty string is returned. If the
// key does exist in the schema, also exists in the map, and that value from the
// map can be coerced to a string, it will be returned. If, however, the key
// does not exist in the map or its value cannot be coerced to a string, any
// default value specified for that field in the schema will be returned. If
// no default is defined, an empty string is returned.
func (p *Parameters) GetString(key string) string {
	if p.Schema == nil {
		return ""
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return ""
		}
	}
	stringSchema, ok := schema.(*StringPropertySchema)
	if !ok {
		return ""
	}
	return ifaceToString(p.Data[key], stringSchema.DefaultValue)
}

// GetStringArray retrieves a []string by key from the Parameters' underlying
// map. If the key does not exist in the schema, nil is returned. If the key
// does exist in the schema, also exists in the map, and that value from the map
// can be coerced to a []string, it will be returned. If, however, the key
// does not exist in the map or its value cannot be coerced to a []string, any
// default value specified for that field in the schema will be returned. If
// no default is defined, nil is returned.
func (p *Parameters) GetStringArray(key string) []string {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return nil
		}
	}
	arrSchema, ok := schema.(*ArrayPropertySchema)
	if !ok {
		return nil
	}
	itemDefault := ""
	if arrSchema.ItemsSchema != nil {
		var itemSchema *StringPropertySchema
		itemSchema, ok = arrSchema.ItemsSchema.(*StringPropertySchema)
		if !ok {
			return nil
		}
		itemDefault = itemSchema.DefaultValue
	}
	valIface, ok := p.Data[key]
	if !ok {
		return ifaceArrayToStringArray(arrSchema.DefaultValue, itemDefault)
	}
	val, ok := valIface.([]interface{})
	if !ok {
		return ifaceArrayToStringArray(arrSchema.DefaultValue, itemDefault)
	}
	return ifaceArrayToStringArray(val, itemDefault)
}

func ifaceArrayToStringArray(arr []interface{}, itemDefault string) []string {
	if len(arr) == 0 {
		return nil
	}
	retArr := make([]string, len(arr))
	for i, item := range arr {
		retArr[i] = ifaceToString(item, itemDefault)
	}
	return retArr
}

func ifaceToString(valIface interface{}, defaultVal string) string {
	if valIface == nil {
		return defaultVal
	}
	if val, ok := valIface.(*string); ok {
		return *val
	}
	if val, ok := valIface.(string); ok {
		return val
	}
	return defaultVal
}

// GetInt64 retrieves an int64 by key from the Parameters' underlying map. If
// the key does not exist in the schema, 0 is returned. If the key does exist in
// the schema, also exists in the map, and that value from the map can be
// coerced to an int64, it will be returned. If, however, the key does not exist
// in the map or its value cannot be coerced to an int64, any default value
// specified for that field in the schema will be returned. If no default is
// defined, 0 is returned.
func (p *Parameters) GetInt64(key string) int64 {
	if p.Schema == nil {
		return 0
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return 0
		}
	}
	intSchema, ok := schema.(*IntPropertySchema)
	if !ok {
		return 0
	}
	return ifaceToInt64(p.Data[key], intSchema.DefaultValue)
}

// GetInt64Array retrieves an []int64 by key from the Parameters' underlying
// map. If the key does not exist in the schema, nil is returned. If the key
// does exist in the schema, also exists in the map, and that value from the map
// can be coerced to an []int64, it will be returned. If, however, the key does
// not exist in the map or its value cannot be coerced to an []int64, any
// default value specified for that field in the schema will be returned. If no
// default is defined, nil is returned.
func (p *Parameters) GetInt64Array(key string) []int64 {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return nil
		}
	}
	arrSchema, ok := schema.(*ArrayPropertySchema)
	if !ok {
		return nil
	}
	var itemDefault *int64
	if arrSchema.ItemsSchema != nil {
		var itemSchema *IntPropertySchema
		itemSchema, ok = arrSchema.ItemsSchema.(*IntPropertySchema)
		if !ok {
			return nil
		}
		itemDefault = itemSchema.DefaultValue
	}
	valIface, ok := p.Data[key]
	if !ok {
		return ifaceArrayToInt64Array(arrSchema.DefaultValue, itemDefault)
	}
	val, ok := valIface.([]interface{})
	if !ok {
		return ifaceArrayToInt64Array(arrSchema.DefaultValue, itemDefault)
	}
	return ifaceArrayToInt64Array(val, itemDefault)
}

func ifaceArrayToInt64Array(arr []interface{}, itemDefault *int64) []int64 {
	if len(arr) == 0 {
		return nil
	}
	retArr := make([]int64, len(arr))
	for i, item := range arr {
		retArr[i] = ifaceToInt64(item, itemDefault)
	}
	return retArr
}

func ifaceToInt64(valIface interface{}, defaultVal *int64) int64 {
	if valIface == nil {
		if defaultVal == nil {
			return 0
		}
		return *defaultVal
	}
	if val, ok := valIface.(*float64); ok {
		return int64(*val)
	}
	if val, ok := valIface.(float64); ok {
		return int64(val)
	}
	if val, ok := valIface.(*float32); ok {
		return int64(*val)
	}
	if val, ok := valIface.(float32); ok {
		return int64(val)
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
	if defaultVal == nil {
		return 0
	}
	return *defaultVal
}

// GetFloat64 retrieves a float64 by key from the Parameters' underlying map. If
// the key does not exist in the schema, 0 is returned. If the key does exist in
// the schema, also exists in the map, and that value from the map can be
// coerced to a float64, it will be returned. If, however, the key does not
// exist in the map or its value cannot be coerced to a float64, any default
// value specified for that field in the schema will be returned. If no default
// is defined, 0 is returned.
func (p *Parameters) GetFloat64(key string) float64 {
	if p.Schema == nil {
		return 0
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return 0
		}
	}
	floatSchema, ok := schema.(*FloatPropertySchema)
	if !ok {
		return 0
	}
	return ifaceToFloat64(p.Data[key], floatSchema.DefaultValue)
}

// GetFloat64Array retrieves a []float64 by key from the Parameters' underlying
// map. If the key does not exist in the schema, nil is returned. If the key
// does exist in the schema, also exists in the map, and that value from the map
// can be coerced to a []float64, it will be returned. If, however, the key does
// not exist in the map or its value cannot be coerced to a []float64, any
// default value specified for that field in the schema will be returned. If no
// default is defined, nil is returned.
func (p *Parameters) GetFloat64Array(key string) []float64 {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return nil
		}
	}
	arrSchema, ok := schema.(*ArrayPropertySchema)
	if !ok {
		return nil
	}
	var itemDefault *float64
	if arrSchema.ItemsSchema != nil {
		var itemSchema *FloatPropertySchema
		itemSchema, ok = arrSchema.ItemsSchema.(*FloatPropertySchema)
		if !ok {
			return nil
		}
		itemDefault = itemSchema.DefaultValue
	}
	valIface, ok := p.Data[key]
	if !ok {
		return ifaceArrayToFloat64Array(arrSchema.DefaultValue, itemDefault)
	}
	val, ok := valIface.([]interface{})
	if !ok {
		return ifaceArrayToFloat64Array(arrSchema.DefaultValue, itemDefault)
	}
	return ifaceArrayToFloat64Array(val, itemDefault)
}

func ifaceArrayToFloat64Array(
	arr []interface{},
	itemDefault *float64,
) []float64 {
	if len(arr) == 0 {
		return nil
	}
	retArr := make([]float64, len(arr))
	for i, item := range arr {
		retArr[i] = ifaceToFloat64(item, itemDefault)
	}
	return retArr
}

func ifaceToFloat64(valIface interface{}, defaultVal *float64) float64 {
	if valIface == nil {
		if defaultVal == nil {
			return 0
		}
		return *defaultVal
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
	if defaultVal == nil {
		return 0
	}
	return *defaultVal
}

// GetObject retrieves a map[string]interface{} by key from the Parameters'
// underlying map wrapped in a new Parameters object. If the key does not exist
// in the schema, a Parameters object with no underlying map is returned. If the
// key does exist in the schema, also exists in the map, and that value from the
// map can be coerced to a map[string]interface{}, it will be wrapped and
// returned. If, however, the key does not exist in the map or its value cannot
// be coerced to a map[string]interface{}, any default value specified for that
// field in the schema will be wrapped and returned. If no default is defined, a
// Parameters object with no underlying map is returned. GetObject calls can
// be chained to "drill down" into a complex set of parameters.
func (p *Parameters) GetObject(key string) Parameters {
	params := Parameters{
		Codec: p.Codec,
	}
	if p.Schema == nil {
		return params
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return params
		}
	}
	objectSchema, ok := schema.(*ObjectPropertySchema)
	if !ok {
		return params
	}
	params.Schema = objectSchema
	return ifaceToParams(
		p.Data[key],
		p.Codec,
		objectSchema,
		objectSchema.DefaultValue,
	)
}

// GetObjectArray retrieves a []map[string]interface{} by key from the
// Parameters' underlying map, wraps each element in a new Parameters object and
// returns []Parameters containing those wrapped []map[string]interface{}. If
// the key does not exist in the schema, nil is returned. If the key does exist
// in the schema, also exists in the map, and that value from the map can be
// coerced to a []map[string]interface{}, each element will be wrapped in a new
// Parameters object and []Parameters will be returned. If, however, the key
// does not exist in the map or its value cannot be coerced to a
// []map[string]interface{}, any default value specified for that
// field in the schema will have each of its elements wrapped and []Parameters
// will be returned. If no default is defined, nil is returned.
func (p *Parameters) GetObjectArray(key string) []Parameters {
	if p.Schema == nil {
		return nil
	}
	schema, ok := p.Schema.GetPropertySchemas()[key]
	if !ok {
		// Schema for the key wasn't found, but maybe "additional" properties are
		// supported?
		schema = p.Schema.GetAdditionalPropertySchema()
		if schema == nil {
			return nil
		}
	}
	arrSchema, ok := schema.(*ArrayPropertySchema)
	if !ok {
		return nil
	}
	var itemDefault map[string]interface{}
	var itemSchema *ObjectPropertySchema
	if arrSchema.ItemsSchema != nil {
		itemSchema, ok = arrSchema.ItemsSchema.(*ObjectPropertySchema)
		if !ok {
			return nil
		}
		itemDefault = itemSchema.DefaultValue
	}
	valIface, ok := p.Data[key]
	if !ok {
		return ifaceArrayToParamsArray(
			arrSchema.DefaultValue,
			p.Codec,
			itemSchema,
			itemDefault,
		)
	}
	val, ok := valIface.([]interface{})
	if !ok {
		return ifaceArrayToParamsArray(
			arrSchema.DefaultValue,
			p.Codec,
			itemSchema,
			itemDefault,
		)
	}
	return ifaceArrayToParamsArray(
		val,
		p.Codec,
		itemSchema,
		itemDefault,
	)
}

func ifaceArrayToParamsArray(
	arr []interface{},
	codec crypto.Codec,
	itemSchema *ObjectPropertySchema, // nolint: interfacer
	itemDefault map[string]interface{},
) []Parameters {
	if len(arr) == 0 {
		return nil
	}
	retArr := make([]Parameters, len(arr))
	for i, item := range arr {
		retArr[i] = ifaceToParams(
			item,
			codec,
			itemSchema,
			itemDefault,
		)
	}
	return retArr
}

func ifaceToParams(
	valIface interface{},
	codec crypto.Codec,
	schema KeyedPropertySchemaContainer,
	defaultVal map[string]interface{},
) Parameters {
	params := Parameters{
		Codec:  codec,
		Schema: schema,
	}
	if valIface == nil {
		params.Data = defaultVal
		return params
	}
	if val, ok := valIface.(map[string]interface{}); ok {
		params.Data = val
		return params
	}
	params.Data = defaultVal
	return params
}
