package types

import (
	"reflect"
)

// IsEmpty determines if the given value is empty or not.
// Uses reflection to check for empty on strings, slice, etc
// We don't want to actually check for zero values on everything
// because that could be a valid case
func IsEmpty(value interface{}) bool {
	if value != nil {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String:
			v := value.(string)
			if v != "" {
				return false

			}
		case reflect.Slice:
			return v.Len() <= 0
		case reflect.Array:
			return v.Len() <= 0
		case reflect.Map:
			return v.Len() <= 0
		default:
			return false
		}
	}
	return true
}
