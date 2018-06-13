package service

import (
	"encoding/json"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"

	"github.com/Azure/open-service-broker-azure/pkg/crypto/fake"
	"github.com/stretchr/testify/assert"
)

func TestMarshalParametersWithMissingSchema(t *testing.T) {
	p := Parameters{
		Codec: fake.NewCodec(),
	}
	_, err := json.Marshal(p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot marshal without a schema")
}

func TestMarshalParametersWithMissingCodec(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{},
	}
	_, err := json.Marshal(p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot marshal without a codec")
}

func TestMarshalParametersWithFielsdNotInSchema(t *testing.T) {
	p := Parameters{
		Codec:  fake.NewCodec(),
		Schema: &InputParametersSchema{},
		Data: map[string]interface{}{
			"foo": "bar",
			"bat": "baz",
		},
	}
	jsonBytes, err := json.Marshal(p)
	assert.Nil(t, err)
	// Convert back to a map to make easier assertions
	mp := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &mp)
	assert.Nil(t, err)
	// There should be nothing in here
	assert.Empty(t, mp)
}

func TestMarshalParametersWithInsecureFields(t *testing.T) {
	codec := fake.NewCodec().(*fake.Codec)
	var encryptCallCount int
	codec.EncryptBehavior = func(plaintext []byte) ([]byte, error) {
		encryptCallCount++
		return plaintext, nil
	}
	p := Parameters{
		Codec: codec,
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
				"bat": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
			"bat": "baz",
		},
	}
	jsonBytes, err := json.Marshal(p)
	assert.Nil(t, err)
	// Convert back to a map to make easier assertions
	mp := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &mp)
	assert.Nil(t, err)
	// There should be exactly two elements
	assert.Equal(t, 2, len(mp))
	// Encrypt should never have been called
	assert.Equal(t, 0, encryptCallCount)
}

func TestMarshalParametersWithNonStringSecureField(t *testing.T) {
	p := Parameters{
		Codec: fake.NewCodec(),
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo"},
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	_, err := json.Marshal(p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot encrypt non-string field")
}

func TestMarshalParametersWithNonStringSecureFieldValue(t *testing.T) {
	p := Parameters{
		Codec: fake.NewCodec(),
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo"},
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	_, err := json.Marshal(p)
	assert.NotNil(t, err)
	assert.Contains(
		t,
		err.Error(),
		"cannot encrypt non-string value of string field",
	)
}

func TestMarshalParametersWithSomeSecureFields(t *testing.T) {
	codec := fake.NewCodec().(*fake.Codec)
	var encryptCallCount int
	codec.EncryptBehavior = func(plaintext []byte) ([]byte, error) {
		encryptCallCount++
		return plaintext, nil
	}
	p := Parameters{
		Codec: codec,
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo", "bat"},
			PropertySchemas: map[string]PropertySchema{
				"abc": &StringPropertySchema{}, // Not secure
				"foo": &StringPropertySchema{},
				"bat": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"abc": "xyz", // Not secure
			"foo": "bar",
			"bat": "baz",
		},
	}
	jsonBytes, err := json.Marshal(p)
	assert.Nil(t, err)
	// Convert back to a map to make easier assertions
	mp := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &mp)
	assert.Nil(t, err)
	// There should be exactly three elements
	assert.Equal(t, 3, len(mp))
	// Encrypt should have been called twice
	assert.Equal(t, 2, encryptCallCount)
}

func TestUnmarshalParametersWithMissingSchema(t *testing.T) {
	p := Parameters{
		Codec: fake.NewCodec(),
	}
	err := json.Unmarshal([]byte("{}"), &p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal without a schema")
}

func TestUnmarshalParametersWithMissingCodec(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{},
	}
	err := json.Unmarshal([]byte("{}"), &p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal without a codec")
}

func TestUnmarshalParametersWithFielsdNotInSchema(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"bat": "baz",
	}
	// Turn the raw map into JSON
	jsonBytes, err := json.Marshal(data)
	assert.Nil(t, err)
	p := Parameters{
		Codec:  fake.NewCodec(),
		Schema: &InputParametersSchema{},
	}
	// Unmarshal into p
	err = json.Unmarshal(jsonBytes, &p)
	assert.Nil(t, err)
	// There should be nothing in p.Data
	assert.Empty(t, p.Data)
}

func TestUnmarshalParametersWithInsecureFields(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
		"bat": "baz",
	}
	// Turn the raw map into JSON
	jsonBytes, err := json.Marshal(data)
	assert.Nil(t, err)
	codec := fake.NewCodec().(*fake.Codec)
	var decryptCallCount int
	codec.DecryptBehavior = func(plaintext []byte) ([]byte, error) {
		decryptCallCount++
		return plaintext, nil
	}
	p := Parameters{
		Codec: codec,
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
				"bat": &StringPropertySchema{},
			},
		},
	}
	err = json.Unmarshal(jsonBytes, &p)
	assert.Nil(t, err)
	// There should be exactly two elements
	assert.Equal(t, 2, len(p.Data))
	// Decrypt should never have been called
	assert.Equal(t, 0, decryptCallCount)
}

func TestUnmarshalParametersWithNonStringSecureField(t *testing.T) {
	data := map[string]interface{}{
		"foo": "bar",
	}
	// Turn the raw map into JSON
	jsonBytes, err := json.Marshal(data)
	assert.Nil(t, err)
	p := Parameters{
		Codec: fake.NewCodec(),
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo"},
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{},
			},
		},
	}
	err = json.Unmarshal(jsonBytes, &p)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot decrypt non-string field")
}

func TestUnmarshalParametersWithNonStringSecureFieldValue(t *testing.T) {
	data := map[string]interface{}{
		"foo": 42,
	}
	// Turn the raw map into JSON
	jsonBytes, err := json.Marshal(data)
	assert.Nil(t, err)
	p := Parameters{
		Codec: fake.NewCodec(),
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo"},
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
			},
		},
	}
	err = json.Unmarshal(jsonBytes, &p)
	assert.NotNil(t, err)
	assert.Contains(
		t,
		err.Error(),
		"cannot decrypt non-string value of string field",
	)
}

func TestUnmarshalParametersWithSomeSecureFields(t *testing.T) {
	data := map[string]interface{}{
		"abc": "xyz", // Not secure
		"foo": "bar",
		"bat": "baz",
	}
	// Turn the raw map into JSON
	jsonBytes, err := json.Marshal(data)
	assert.Nil(t, err)
	codec := fake.NewCodec().(*fake.Codec)
	var dectypeCallCount int
	codec.DecryptBehavior = func(plaintext []byte) ([]byte, error) {
		dectypeCallCount++
		return plaintext, nil
	}
	p := Parameters{
		Codec: codec,
		Schema: &InputParametersSchema{
			SecureProperties: []string{"foo", "bat"},
			PropertySchemas: map[string]PropertySchema{
				"abc": &StringPropertySchema{}, // Not secure
				"foo": &StringPropertySchema{},
				"bat": &StringPropertySchema{},
			},
		},
	}
	err = json.Unmarshal(jsonBytes, &p)
	assert.Nil(t, err)
	// There should be exactly three elements
	assert.Equal(t, 3, len(p.Data))
	// Encrypt should have been called twice
	assert.Equal(t, 2, dectypeCallCount)
}

func TestGetStringWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
}

func TestGetStringWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
}

func TestGetStringWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional:      &StringPropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "bar", val)
}

func TestGetStringWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
			},
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
}

func TestGetStringNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
}

func TestGetStringNotInMapWithDefault(t *testing.T) {
	const defaultVal = "bar"
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{
					DefaultValue: defaultVal,
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetString("foo")
	assert.Equal(t, defaultVal, val)
}

func TestGetStringValueIsNotString(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
}

func TestGetString(t *testing.T) {
	const expectedVal = "bar"
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo1": &StringPropertySchema{},
				"foo2": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo1": ptr.ToString(expectedVal),
			"foo2": expectedVal,
		},
	}
	val := p.GetString("foo1")
	assert.Equal(t, expectedVal, val)
	val = p.GetString("foo2")
	assert.Equal(t, expectedVal, val)
}

func TestGetStringArrayWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": []interface{}{"bat", "baz"},
		},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArrayWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{"bat", "baz"},
		},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArrayWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional: &ArrayPropertySchema{
				ItemsSchema: &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{"bat", "baz"},
		},
	}
	val := p.GetStringArray("foo")
	assert.Equal(t, []string{"bat", "baz"}, val)
}

func TestGetStringArrayWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArrayNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArrayNotInMapWithDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					DefaultValue: []interface{}{"bar", "bat"},
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetStringArray("foo")
	assert.Equal(t, []string{"bar", "bat"}, val)
}

func TestGetStringArrayValueIsNotArray(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArrayWithNonStringsInArray(t *testing.T) {
	const defaultString = "bar"
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: &StringPropertySchema{
						DefaultValue: defaultString,
					},
				},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{42},
		},
	}
	val := p.GetStringArray("foo")
	assert.Equal(t, []string{defaultString}, val)
}

func TestGetStringArray(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{
				ptr.ToString("bar"),
				"bat",
			},
		},
	}
	val := p.GetStringArray("foo")
	assert.Equal(t, []string{"bar", "bat"}, val)
}

func TestGetInt64WithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
}

func TestGetInt64WithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
}

func TestGetInt64WithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	const expectedValue int64 = 42
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional:      &IntPropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": expectedValue,
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, expectedValue, val)
}

func TestGetInt64WithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{},
			},
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
}

func TestGetInt64NotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
}

func TestGetInt64NotInMapWithDefault(t *testing.T) {
	const defaultVal int64 = 42
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{
					DefaultValue: ptr.ToInt64(defaultVal),
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, defaultVal, val)
}

func TestGetInt64ValueIsNotInt(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &IntPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
}

func TestGetInt64(t *testing.T) {
	const expectedVal int64 = 42
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo1": &IntPropertySchema{},
				"foo2": &IntPropertySchema{},
				"foo3": &IntPropertySchema{},
				"foo4": &IntPropertySchema{},
				"foo5": &IntPropertySchema{},
				"foo6": &IntPropertySchema{},
			},
		},
		// All of the following cases should be tolerated
		Data: map[string]interface{}{
			"foo1": ptr.ToInt64(expectedVal),
			"foo2": expectedVal,
			"foo3": ptr.ToInt32(int32(expectedVal)),
			"foo4": int32(expectedVal),
			"foo5": ptr.ToInt(int(expectedVal)),
			"foo6": int(expectedVal),
		},
	}
	val := p.GetInt64("foo1")
	assert.Equal(t, expectedVal, val)
	val = p.GetInt64("foo2")
	assert.Equal(t, expectedVal, val)
	val = p.GetInt64("foo3")
	assert.Equal(t, expectedVal, val)
	val = p.GetInt64("foo4")
	assert.Equal(t, expectedVal, val)
	val = p.GetInt64("foo5")
	assert.Equal(t, expectedVal, val)
	val = p.GetInt64("foo6")
	assert.Equal(t, expectedVal, val)
}

func TestGetInt64ArrayWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": []interface{}{8, 42},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Nil(t, val)
}

func TestGetInt64ArrayWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{8, 42},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Nil(t, val)
}

func TestGetInt64ArrayWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional: &ArrayPropertySchema{
				ItemsSchema: &IntPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{8, 42},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Equal(t, []int64{8, 42}, val)
}

func TestGetInt64ArrayWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Nil(t, val)
}

func TestGetInt64ArrayNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetInt64Array("foo")
	assert.Nil(t, val)
}

func TestGetInt64ArrayNotInMapWithDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					DefaultValue: []interface{}{8, 42},
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetInt64Array("foo")
	assert.Equal(t, []int64{8, 42}, val)
}

func TestGetInt64ArrayValueIsNotArray(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	val := p.GetInt64Array("foo")
	assert.Nil(t, val)
}

func TestGetInt64ArrayWithNonIntsInArray(t *testing.T) {
	const defaultInt int64 = 42
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: &IntPropertySchema{
						DefaultValue: ptr.ToInt64(defaultInt),
					},
				},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{"bar"},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Equal(t, []int64{defaultInt}, val)
}

func TestGetInt64Array(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			// All of these should work...
			"foo": []interface{}{
				ptr.ToInt64(1),
				int64(2),
				ptr.ToInt32(3),
				int32(4),
				ptr.ToInt(5),
				6,
			},
		},
	}
	val := p.GetInt64Array("foo")
	assert.Equal(
		t,
		[]int64{1, 2, 3, 4, 5, 6},
		val,
	)
}

func TestGetFloat64WithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": 3.14,
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
}

func TestGetFloat64WithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": 3.14,
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
}

func TestGetFloat64WithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	const expectedValue float64 = 3.14
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional:      &FloatPropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": expectedValue,
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, expectedValue, val)
}

func TestGetFloat64WithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &FloatPropertySchema{},
			},
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
}

func TestGetFloat64NotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &FloatPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
}

func TestGetFloat64NotInMapWithDefault(t *testing.T) {
	const defaultVal float64 = 3.14
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &FloatPropertySchema{
					DefaultValue: ptr.ToFloat64(defaultVal),
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, defaultVal, val)
}

func TestGetInt64ValueIsNotFloat(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &FloatPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
}

func TestGetFloat64(t *testing.T) {
	const expectedFooVal float32 = 3.14
	const expectedBarVal int64 = 42
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo1": &FloatPropertySchema{},
				"foo2": &FloatPropertySchema{},
				"foo3": &FloatPropertySchema{},
				"foo4": &FloatPropertySchema{},
				"bar1": &FloatPropertySchema{},
				"bar2": &FloatPropertySchema{},
				"bar3": &FloatPropertySchema{},
				"bar4": &FloatPropertySchema{},
				"bar5": &FloatPropertySchema{},
				"bar6": &FloatPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			// All of the following cases should be tolerated
			"foo1": ptr.ToFloat64(float64(expectedFooVal)),
			"foo2": float64(expectedFooVal),
			"foo3": ptr.ToFloat32(expectedFooVal),
			"foo4": expectedFooVal,
			// Mathematically, ints are floats too, ya know...
			"bar1": ptr.ToInt64(expectedBarVal),
			"bar2": expectedBarVal,
			"bar3": ptr.ToInt32(int32(expectedBarVal)),
			"bar4": int32(expectedBarVal),
			"bar5": ptr.ToInt(int(expectedBarVal)),
			"bar6": int(expectedBarVal),
		},
	}
	val := p.GetFloat64("foo1")
	assert.Equal(t, float64(expectedFooVal), val)
	val = p.GetFloat64("foo2")
	assert.Equal(t, float64(expectedFooVal), val)
	val = p.GetFloat64("foo3")
	assert.Equal(t, float64(expectedFooVal), val)
	val = p.GetFloat64("foo4")
	assert.Equal(t, float64(expectedFooVal), val)
	val = p.GetFloat64("bar1")
	assert.Equal(t, float64(expectedBarVal), val)
	val = p.GetFloat64("bar2")
	assert.Equal(t, float64(expectedBarVal), val)
	val = p.GetFloat64("bar3")
	assert.Equal(t, float64(expectedBarVal), val)
	val = p.GetFloat64("bar4")
	assert.Equal(t, float64(expectedBarVal), val)
	val = p.GetFloat64("bar5")
	assert.Equal(t, float64(expectedBarVal), val)
	val = p.GetFloat64("bar6")
	assert.Equal(t, float64(expectedBarVal), val)
}

func TestGetFloat64ArrayWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": []interface{}{3.14, 8, 42},
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Nil(t, val)
}

func TestGetFloat64ArrayWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) { // nolint: lll
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{3.14, 8, 42},
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Nil(t, val)
}

func TestGetFloat64ArrayWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional: &ArrayPropertySchema{
				ItemsSchema: &FloatPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{3.14, 8, 42},
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Equal(t, []float64{3.14, 8, 42}, val)
}

func TestGetFloat64ArrayWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Nil(t, val)
}

func TestGetFloat64ArrayNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetFloat64Array("foo")
	assert.Nil(t, val)
}

func TestGetFloat64ArrayNotInMapWithDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					DefaultValue: []interface{}{3.14, 8, 42},
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetFloat64Array("foo")
	assert.Equal(t, []float64{3.14, 8, 42}, val)
}

func TestGetFloat64ArrayValueIsNotArray(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": 3.14,
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Nil(t, val)
}

func TestGetFloat64Array(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			// All of these should work...
			"foo": []interface{}{
				ptr.ToFloat64(3.14),
				float64(3.14),
				// These can't be tested nicely because of floating point errors when
				// the function under test converts float32 to float64. I (krancour)
				// am going to let this slide for now. afaik, in "real life," maps
				// representing param data will always contain float64s due to how
				// Go unmarshals JSON into maps.
				// ptr.ToFloat32(3.14),
				// float32(3.14),
				ptr.ToInt64(1),
				int64(2),
				ptr.ToInt32(3),
				int32(4),
				ptr.ToInt(5),
				6,
			},
		},
	}
	val := p.GetFloat64Array("foo")
	assert.Equal(
		t,
		[]float64{3.14, 3.14, 1, 2, 3, 4, 5, 6},
		val,
	)
}

func TestGetObjectWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	}
	val := p.GetObject("foo")
	assert.Equal(t, Parameters{}, val)
}

func TestGetObjectWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	}
	val := p.GetObject("foo")
	assert.Equal(t, Parameters{}, val)
}

func TestGetObjectWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	expectedVal := map[string]interface{}{}
	additionalPropertySchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional:      additionalPropertySchema,
		},
		Data: map[string]interface{}{
			"foo": expectedVal,
		},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: additionalPropertySchema,
			Data:   expectedVal,
		},
		val,
	)
}

func TestGetObjectWithNoMap(t *testing.T) {
	fooSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": fooSchema,
			},
		},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: fooSchema,
		},
		val,
	)
}

func TestGetObjectNotInMapWithNoDefault(t *testing.T) {
	fooSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": fooSchema,
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: fooSchema,
		},
		val,
	)
}

func TestGetObjectNotInMapWithDefault(t *testing.T) {
	defaultVal := map[string]interface{}{
		"bar": "bat",
	}
	fooSchema := &ObjectPropertySchema{
		DefaultValue: defaultVal,
	}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": fooSchema,
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: fooSchema,
			Data:   defaultVal,
		},
		val,
	)
}

func TestGetObjectValueIsNotMap(t *testing.T) {
	fooSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": fooSchema,
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: fooSchema,
		},
		val,
	)
}

func TestGetObject(t *testing.T) {
	expectedVal := map[string]interface{}{
		"bar": "bat",
	}
	fooSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": fooSchema,
			},
		},
		Data: map[string]interface{}{
			"foo": expectedVal,
		},
	}
	val := p.GetObject("foo")
	assert.Equal(
		t,
		Parameters{
			Schema: fooSchema,
			Data:   expectedVal,
		},
		val,
	)
}

func TestGetObjectArrayWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Nil(t, val)
}

func TestGetObjectArrayWithNoSchemaForKeyAndAdditionalNotAllowed(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{
				map[string]interface{}{},
			},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Nil(t, val)
}

func TestGetObjectArrayWithNoSchemaForKeyAndAdditionalAllowed(t *testing.T) {
	additionalItemSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &ObjectPropertySchema{
			PropertySchemas: map[string]PropertySchema{},
			Additional: &ArrayPropertySchema{
				ItemsSchema: additionalItemSchema,
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{
				map[string]interface{}{
					"bar": "bat",
				},
				map[string]interface{}{
					"bat": "baz",
				},
			},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Equal(
		t,
		[]Parameters{
			{
				Schema: additionalItemSchema,
				Data: map[string]interface{}{
					"bar": "bat",
				},
			},
			{
				Schema: additionalItemSchema,
				Data: map[string]interface{}{
					"bat": "baz",
				},
			},
		},
		val,
	)
}

func TestGetObjectArrayWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Nil(t, val)
}

func TestGetObjectArrayNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetObjectArray("foo")
	assert.Nil(t, val)
}

func TestGetObjectArrayNotInMapWithDefault(t *testing.T) {
	fooItemSchema := &ObjectPropertySchema{
		PropertySchemas: map[string]PropertySchema{
			"bar": &StringPropertySchema{},
		},
	}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: fooItemSchema,
					DefaultValue: []interface{}{
						map[string]interface{}{
							"bar": "bat",
						},
						map[string]interface{}{
							"bar": "baz",
						},
					},
				},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetObjectArray("foo")
	assert.Equal(
		t,
		[]Parameters{
			{
				Schema: fooItemSchema,
				Data: map[string]interface{}{
					"bar": "bat",
				},
			},
			{
				Schema: fooItemSchema,
				Data: map[string]interface{}{
					"bar": "baz",
				},
			},
		},
		val,
	)
}

func TestGetObjectArrayValueIsNotArray(t *testing.T) {
	fooItemSchema := &ObjectPropertySchema{
		PropertySchemas: map[string]PropertySchema{
			"bar": &StringPropertySchema{},
		},
	}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: fooItemSchema,
				},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetObjectArray("foo")
	assert.Nil(t, val)
}

func TestGetObjectArrayWithNonMapsInArray(t *testing.T) {
	defaultObject := map[string]interface{}{
		"foo": "bar",
		"bat": "baz",
	}
	fooItemSchema := &ObjectPropertySchema{
		DefaultValue: defaultObject,
	}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: fooItemSchema,
				},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{"bar"},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Equal(
		t,
		[]Parameters{
			{
				Schema: fooItemSchema,
				Data:   defaultObject,
			},
		},
		val,
	)
}

func TestGetObjectArray(t *testing.T) {
	foo1Value := map[string]interface{}{
		"foo": "bar",
	}
	foo2Value := map[string]interface{}{
		"foo": "bat",
	}
	fooItemSchema := &ObjectPropertySchema{}
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{
					ItemsSchema: fooItemSchema,
				},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{
				foo1Value,
				foo2Value,
			},
		},
	}
	val := p.GetObjectArray("foo")
	assert.Equal(
		t,
		[]Parameters{
			{
				Schema: fooItemSchema,
				Data:   foo1Value,
			},
			{
				Schema: fooItemSchema,
				Data:   foo2Value,
			},
		},
		val,
	)
}
