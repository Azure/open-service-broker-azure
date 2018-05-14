package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/stretchr/testify/assert"
)

func TestStringPropertySchemaToJSON(t *testing.T) {
	fooSps := StringPropertySchema{
		Description:   "foo",
		AllowedValues: []string{"foo", "bar", "bat", "baz"},
	}
	jsonBytes, err := json.Marshal(fooSps)
	assert.Nil(t, err)
	// We'll unmarshal into a map (that should always work) and then we'll
	// make assertions on the map to prove the JSON was what we'd expected it to
	// be.
	fooSpsMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &fooSpsMap)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(fooSpsMap))
	schemaType, ok := fooSpsMap["type"]
	assert.True(t, ok)
	assert.Equal(t, "string", schemaType)
	description, ok := fooSpsMap["description"]
	assert.True(t, ok)
	assert.Equal(t, fooSps.Description, description)
	allowedValues, ok := fooSpsMap["enum"]
	assert.True(t, ok)
	allowedValuesIfaces, ok := allowedValues.([]interface{})
	assert.True(t, ok)
	allowedValuesStrings := make([]string, len(allowedValuesIfaces))
	assert.Equal(t, len(fooSps.AllowedValues), len(allowedValuesStrings))
	for i, allowedValueIface := range allowedValuesIfaces {
		var ok bool
		allowedValuesStrings[i], ok = allowedValueIface.(string)
		assert.True(t, ok)
	}
	assert.Equal(t, fooSps.AllowedValues, allowedValuesStrings)
}

func TestValidateStringProperty(t *testing.T) {
	const fieldName = "xyz"

	sps := StringPropertySchema{}
	err := sps.validate(fieldName, 5)
	assert.NotNil(t, err)
	validationError, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	sps = StringPropertySchema{
		MinLength: ptr.ToInt(3),
	}
	err = sps.validate(fieldName, "fo")
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = sps.validate(fieldName, "foo")
	assert.Nil(t, err)
	err = sps.validate(fieldName, "foobar")
	assert.Nil(t, err)

	sps = StringPropertySchema{
		MaxLength: ptr.ToInt(6),
	}
	err = sps.validate(fieldName, "foobarr")
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = sps.validate(fieldName, "foobar")
	assert.Nil(t, err)
	err = sps.validate(fieldName, "foo")
	assert.Nil(t, err)

	sps = StringPropertySchema{
		AllowedValues: []string{"foo", "bar"},
	}
	err = sps.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = sps.validate(fieldName, "foo")
	assert.Nil(t, err)
	err = sps.validate(fieldName, "bar")
	assert.Nil(t, err)

	sps = StringPropertySchema{
		AllowedPattern: regexp.MustCompile(`^\w{3}$`),
	}
	err = sps.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = sps.validate(fieldName, "foo")
	assert.Nil(t, err)
	err = sps.validate(fieldName, "bar")
	assert.Nil(t, err)
}

func TestIntPropertySchemaToJSON(t *testing.T) {
	smallEvenIps := IntPropertySchema{
		Description:      "small, even integers",
		MinValue:         ptr.ToInt64(0),
		MaxValue:         ptr.ToInt64(8),
		AllowedIncrement: ptr.ToInt64(2),
	}
	jsonBytes, err := json.Marshal(smallEvenIps)
	assert.Nil(t, err)
	// We'll unmarshal into a map (that should always work) and then we'll
	// make assertions on the map to prove the JSON was what we'd expected it to
	// be.
	smallEvenIpsMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &smallEvenIpsMap)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(smallEvenIpsMap))
	schemaType, ok := smallEvenIpsMap["type"]
	assert.True(t, ok)
	assert.Equal(t, "integer", schemaType)
	description, ok := smallEvenIpsMap["description"]
	assert.True(t, ok)
	assert.Equal(t, smallEvenIps.Description, description)
	minValueIface, ok := smallEvenIpsMap["minimum"]
	assert.True(t, ok)
	minValueFloat, ok := minValueIface.(float64)
	assert.True(t, ok)
	assert.Equal(t, *smallEvenIps.MinValue, int64(minValueFloat))
	maxValueIface, ok := smallEvenIpsMap["maximum"]
	assert.True(t, ok)
	maxValueFloat, ok := maxValueIface.(float64)
	assert.True(t, ok)
	assert.Equal(t, *smallEvenIps.MaxValue, int64(maxValueFloat))
	incrementIface, ok := smallEvenIpsMap["multipleOf"]
	assert.True(t, ok)
	incrementFloat, ok := incrementIface.(float64)
	assert.True(t, ok)
	assert.Equal(t, *smallEvenIps.AllowedIncrement, int64(incrementFloat))
}

func TestValidateIntProperty(t *testing.T) {
	const fieldName = "xyz"

	ips := IntPropertySchema{}
	err := ips.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	ips = IntPropertySchema{}
	err = ips.validate(fieldName, 3.14)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	ips = IntPropertySchema{
		MinValue: ptr.ToInt64(3),
	}
	err = ips.validate(fieldName, 2.0)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = ips.validate(fieldName, 3.0)
	assert.Nil(t, err)
	err = ips.validate(fieldName, 6.0)
	assert.Nil(t, err)

	ips = IntPropertySchema{
		MaxValue: ptr.ToInt64(6),
	}
	err = ips.validate(fieldName, 7.0)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = ips.validate(fieldName, 6.0)
	assert.Nil(t, err)
	err = ips.validate(fieldName, 3.0)
	assert.Nil(t, err)

	ips = IntPropertySchema{
		AllowedValues: []int64{3, 4},
	}
	err = ips.validate(fieldName, 5.0)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = ips.validate(fieldName, 3.0)
	assert.Nil(t, err)
	err = ips.validate(fieldName, 4.0)
	assert.Nil(t, err)

	ips = IntPropertySchema{
		AllowedIncrement: ptr.ToInt64(2),
	}
	err = ips.validate(fieldName, 5.0)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = ips.validate(fieldName, 0.0)
	assert.Nil(t, err)
	err = ips.validate(fieldName, 8.0)
	assert.Nil(t, err)
}

func TestFloatPropertySchemaToJSON(t *testing.T) {
	smallIntsAndHalvesFps := FloatPropertySchema{
		Description: "small, integers and halves",
		MinValue:    ptr.ToFloat64(0),
		MaxValue:    ptr.ToFloat64(8),
		// AllowedIncrement: ptr.ToFloat64(0.5),
	}
	jsonBytes, err := json.Marshal(smallIntsAndHalvesFps)
	assert.Nil(t, err)
	// We'll unmarshal into a map (that should always work) and then we'll
	// make assertions on the map to prove the JSON was what we'd expected it to
	// be.
	smallEvenIpsMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &smallEvenIpsMap)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(smallEvenIpsMap))
	schemaType, ok := smallEvenIpsMap["type"]
	assert.True(t, ok)
	assert.Equal(t, "number", schemaType)
	description, ok := smallEvenIpsMap["description"]
	assert.True(t, ok)
	assert.Equal(t, smallIntsAndHalvesFps.Description, description)
	minValueIface, ok := smallEvenIpsMap["minimum"]
	assert.True(t, ok)
	minValueFloat, ok := minValueIface.(float64)
	assert.True(t, ok)
	assert.Equal(t, *smallIntsAndHalvesFps.MinValue, float64(minValueFloat))
	maxValueIface, ok := smallEvenIpsMap["maximum"]
	assert.True(t, ok)
	maxValueFloat, ok := maxValueIface.(float64)
	assert.True(t, ok)
	assert.Equal(t, *smallIntsAndHalvesFps.MaxValue, float64(maxValueFloat))
	// incrementIface, ok := smallEvenIpsMap["multipleOf"]
	// assert.True(t, ok)
	// incrementFloat, ok := incrementIface.(float64)
	// assert.True(t, ok)
	// assert.Equal(
	// 	t,
	// 	*smallIntsAndHalvesFps.AllowedIncrement,
	// 	float64(incrementFloat),
	// )
}

func TestValidateFloatProperty(t *testing.T) {
	const fieldName = "xyz"

	fps := FloatPropertySchema{}
	err := fps.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	fps = FloatPropertySchema{
		MinValue: ptr.ToFloat64(3.14),
	}
	err = fps.validate(fieldName, 2.5)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = fps.validate(fieldName, 3.5)
	assert.Nil(t, err)
	err = fps.validate(fieldName, 4.5)
	assert.Nil(t, err)

	fps = FloatPropertySchema{
		MaxValue: ptr.ToFloat64(3.14),
	}
	err = fps.validate(fieldName, 3.5)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = fps.validate(fieldName, 3.0)
	assert.Nil(t, err)
	err = fps.validate(fieldName, 2.5)
	assert.Nil(t, err)

	fps = FloatPropertySchema{
		AllowedValues: []float64{3.14, 4.5},
	}
	err = fps.validate(fieldName, 5.0)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = fps.validate(fieldName, 3.14)
	assert.Nil(t, err)
	err = fps.validate(fieldName, 4.5)
	assert.Nil(t, err)

	fps = FloatPropertySchema{
		// AllowedIncrement: ptr.ToFloat64(0.5),
	}
	// err = fps.validate(fieldName, 4.25)
	// assert.NotNil(t, err)
	// validationError, ok = err.(*ValidationError)
	// assert.True(t, ok)
	// assert.Equal(t, fieldName, validationError.Field)
	err = fps.validate(fieldName, 0.0)
	assert.Nil(t, err)
	err = fps.validate(fieldName, 8.5)
	assert.Nil(t, err)
}

func TestArrayPropertySchemaToJSON(t *testing.T) {
	fooAps := ArrayPropertySchema{
		Description: "a handful of foo",
		ItemsSchema: StringPropertySchema{
			Description:   "foo",
			AllowedValues: []string{"foo", "bar", "bat", "baz"},
		},
	}
	jsonBytes, err := json.Marshal(fooAps)
	assert.Nil(t, err)
	// We'll unmarshal into a map (that should always work) and then we'll
	// make assertions on the map to prove the JSON was what we'd expected it to
	// be.
	fooApsMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &fooApsMap)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(fooApsMap))
	schemaType, ok := fooApsMap["type"]
	assert.True(t, ok)
	assert.Equal(t, "array", schemaType)
	description, ok := fooApsMap["description"]
	assert.True(t, ok)
	assert.Equal(t, fooAps.Description, description)
	_, ok = fooApsMap["items"]
	assert.True(t, ok)
	// We've separately tested StringPropertySchemas, so we won't bother making
	// assertions on "items"
}

func TestValidateArrayProperty(t *testing.T) {
	const fieldName = "xyz"

	aps := ArrayPropertySchema{}
	err := aps.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	aps = ArrayPropertySchema{
		MinItems: ptr.ToInt(3),
	}
	err = aps.validate(fieldName, []interface{}{1, 2})
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = aps.validate(fieldName, []interface{}{1, 2, 3})
	assert.Nil(t, err)
	err = aps.validate(fieldName, []interface{}{1, 2, 3, 4})
	assert.Nil(t, err)

	aps = ArrayPropertySchema{
		MaxItems: ptr.ToInt(6),
	}
	err = aps.validate(fieldName, []interface{}{1, 2, 3, 4, 5, 6, 7})
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)
	err = aps.validate(fieldName, []interface{}{1, 2, 3, 4, 5, 6})
	assert.Nil(t, err)
	err = aps.validate(fieldName, []interface{}{1, 2, 3})
	assert.Nil(t, err)

	aps = ArrayPropertySchema{
		ItemsSchema: IntPropertySchema{
			MinValue: ptr.ToInt64(3),
		},
	}
	err = aps.validate(fieldName, []interface{}{3.0, 2.0, 1.0})
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%s[1]", fieldName), validationError.Field)
	err = aps.validate(fieldName, []interface{}{3.0, 4.0, 5.0})
	assert.Nil(t, err)
}

func TestObjectPropertySchemaToJSON(t *testing.T) {
	myOps := ObjectPropertySchema{
		Description: "a small even integer and a foo",
		PropertySchemas: map[string]PropertySchema{
			"foo": StringPropertySchema{
				Description:   "foo",
				AllowedValues: []string{"foo", "bar", "bat", "baz"},
			},
			"smallEvenInt": IntPropertySchema{
				Description:      "small, even integers",
				MinValue:         ptr.ToInt64(0),
				MaxValue:         ptr.ToInt64(8),
				AllowedIncrement: ptr.ToInt64(2),
			},
		},
	}
	jsonBytes, err := json.Marshal(myOps)
	assert.Nil(t, err)
	// We'll unmarshal into a map (that should always work) and then we'll
	// make assertions on the map to prove the JSON was what we'd expected it to
	// be.
	myOpsMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &myOpsMap)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(myOpsMap))
	schemaType, ok := myOpsMap["type"]
	assert.True(t, ok)
	assert.Equal(t, "object", schemaType)
	description, ok := myOpsMap["description"]
	assert.True(t, ok)
	assert.Equal(t, myOps.Description, description)
	propertiesIface, ok := myOpsMap["properties"]
	assert.True(t, ok)
	propertiesMap, ok := propertiesIface.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(myOps.PropertySchemas), len(propertiesMap))
	// We've separately tested StringPropertySchemas and IntPropertySchemas, so we
	// won't bother making assertions on individual "properties"
}

func TestValidateObjectProperty(t *testing.T) {
	const fieldName = "xyz"

	ops := ObjectPropertySchema{}
	err := ops.validate(fieldName, "foobar")
	assert.NotNil(t, err)
	validationError, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fieldName, validationError.Field)

	ops = ObjectPropertySchema{
		RequiredProperties: []string{"foo", "bat"},
	}
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"foo": "bar",
		},
	)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%s.bat", fieldName), validationError.Field)
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"bat": "baz",
		},
	)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%s.foo", fieldName), validationError.Field)
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"foo": "bar",
			"bat": "baz",
		},
	)
	assert.Nil(t, err)

	ops = ObjectPropertySchema{
		PropertySchemas: map[string]PropertySchema{
			"foo": StringPropertySchema{
				AllowedValues: []string{"bar", "bat", "baz"},
			},
			"bar": IntPropertySchema{
				AllowedIncrement: ptr.ToInt64(2),
			},
		},
	}
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"foo": "bogus",
			"bar": 4.0,
		},
	)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%s.foo", fieldName), validationError.Field)
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"foo": "bar",
			"bar": 5.0,
		},
	)
	assert.NotNil(t, err)
	validationError, ok = err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("%s.bar", fieldName), validationError.Field)
	err = ops.validate(
		fieldName,
		map[string]interface{}{
			"foo": "bar",
			"bar": 4.0,
		},
	)
	assert.Nil(t, err)
}
