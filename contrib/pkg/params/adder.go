package params

import (
	"fmt"
	"strconv"
)

// Add adds the given key/value pair to the given map[string]interface{}, but
// complex keys that specify an index (for multi-valued / array parameters) have
// their values added to that map as a []interface{} under that key. That slice
// is re-sized as needed.
func Add(params map[string]interface{}, key string, val interface{}) error {
	if matches := arrayParamKeyRegex.FindStringSubmatch(key); len(matches) > 0 {
		// Key is indexed
		key = matches[1]
		indexStr := matches[2]
		// Can't get an error here because we already know indexStr is parsable as
		// an integer
		index, _ := strconv.Atoi(indexStr)
		var slice []interface{}
		existingVal, ok := params[key]
		if ok { // If something already exists under this key...
			slice, ok = existingVal.([]interface{})
			if !ok { // And it's NOT a []interface{}...
				// We have a problem
				return fmt.Errorf(
					`key "%s" used in both single-valued and multi-valued context`,
					key,
				)
			}
			// Something already exists here AND it's a slice!
			if len(slice) < index+1 { // If the slice isn't big enough...
				// Grow it!
				// By how much?
				diff := (index + 1) - len(slice)
				slice = append(slice, make([]interface{}, diff)...)
			}
		} else { // If nothing already exists under this key...
			// Make a slice that's big enough
			slice = make([]interface{}, index+1)
		}
		slice[index] = val
		params[key] = slice
	} else {
		// Key isn't indexed
		existingVal, ok := params[key]
		if ok { // If something already exists under this key...
			_, ok := existingVal.([]interface{})
			if ok { // And it's a []interface{}...
				// We have a problem
				return fmt.Errorf(
					`key "%s" used in both single-valued and multi-valued context`,
					key,
				)
			}
		}
		params[key] = val
	}
	return nil
}
