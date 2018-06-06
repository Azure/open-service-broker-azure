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

// -----------------------------------------------------------------------------

func TestGetStringWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, "", val)
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
				"foo": &StringPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	val := p.GetString("foo")
	assert.Equal(t, expectedVal, val)
}

// -----------------------------------------------------------------------------

func TestGetStringArrayWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": []interface{}{"bat", "baz"},
		},
	}
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
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

func TestGetStringArrayValueIsNotStringArray(t *testing.T) {
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
	val := p.GetStringArray("foo")
	assert.Nil(t, val)
}

func TestGetStringArray(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ArrayPropertySchema{},
			},
		},
		Data: map[string]interface{}{
			"foo": []interface{}{"bar", "bat"},
		},
	}
	val := p.GetStringArray("foo")
	assert.Equal(t, []string{"bar", "bat"}, val)
}

// -----------------------------------------------------------------------------

func TestGetInt64WithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": 42,
		},
	}
	val := p.GetInt64("foo")
	assert.Equal(t, int64(0), val)
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

// -----------------------------------------------------------------------------

func TestGetFloat64WithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": 3.14,
		},
	}
	val := p.GetFloat64("foo")
	assert.Equal(t, float64(0), val)
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

// -----------------------------------------------------------------------------

func TestGetObjectWithNoSchema(t *testing.T) {
	p := Parameters{
		Data: map[string]interface{}{
			"foo": map[string]interface{}{},
		},
	}
	val := p.GetObject("foo")
	assert.Equal(t, Parameters{}, val)
}

func TestGetObjectWithNoMap(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ObjectPropertySchema{},
			},
		},
	}
	val := p.GetObject("foo")
	assert.Equal(t, Parameters{}, val)
}

func TestGetObjectNotInMapWithNoDefault(t *testing.T) {
	p := Parameters{
		Schema: &InputParametersSchema{
			PropertySchemas: map[string]PropertySchema{
				"foo": &ObjectPropertySchema{},
			},
		},
		Data: map[string]interface{}{},
	}
	val := p.GetObject("foo")
	assert.Equal(t, Parameters{}, val)
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

// -----------------------------------------------------------------------------
