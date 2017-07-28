package service

import "encoding/json"

// BindingResponse represents the response to a binding request
type BindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

// NewBindingResponseFromJSONString returns a new BindingResponse unmarshalled
// from the provided JSON string
func GetBindingResponseFromJSONString(
	jsonStr string,
	bindingResponse *BindingResponse,
) error {
	return json.Unmarshal([]byte(jsonStr), bindingResponse)
}

// ToJSONString returns a string containing a JSON representation of the
// binding response
func (b *BindingResponse) ToJSONString() (string, error) {
	bytes, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
