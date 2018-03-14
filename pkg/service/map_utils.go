package service

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// GetStructFromMap is a utility function that extracts values from a
// map[string]interface{} to the provided struct. Maps form the input and
// output to all ServiceManager functions. This utility function offers
// ServiceManager implementors the option of working with firendlier,
// service-specific structs by easily extracting information from maps into
// thos structs.
func GetStructFromMap(m map[string]interface{}, s interface{}) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			TagName: "json",
			Result:  s,
		},
	)
	if err != nil {
		return fmt.Errorf("error building tag map decoder: %s", err)
	}
	err = decoder.Decode(m)
	if err != nil {
		return fmt.Errorf("error extracting map to struct: %s", err)
	}
	return nil
}

// GetMapFromStruct is a utility function that extracts values from the provided
// struct into a map[string]interface{}. Maps form the input and output to all
// ServiceManager functions. This utility function offers ServiceManager
// implementors the option of working with firendlier, service-specific structs
// then easily extracting information from those structs into maps.
func GetMapFromStruct(s interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("error getting map from struct: %s", err)
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return nil, fmt.Errorf("error getting map from struct: %s", err)
	}
	return m, nil
}
