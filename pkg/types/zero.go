package types

import (
	"reflect"
)

// IsZero determines if the given value is a zero value
func IsZero(value interface{}) bool {
	if value != nil {
		switch value.(type) {
		case string:
			v := value.(string)
			if v != "" {
				return false

			}
		default:
			zero := reflect.Zero(reflect.TypeOf(value)).Interface()
			isZero := reflect.DeepEqual(value, zero)
			if !isZero {
				return false
			}
		}
	}
	return true
}
