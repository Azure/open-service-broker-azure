package service

import (
	"encoding/json"
	"testing"

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

// ----------------------------------------------------------------------------

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
