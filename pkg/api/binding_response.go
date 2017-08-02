package api

import "encoding/json"

// BindingResponse represents the response to a binding request
type BindingResponse struct {
	Credentials interface{} `json:"credentials"`
}

// GetBindingResponseFromJSON returns a new BindingResponse unmarshalled from
// the provided JSON []byte
func GetBindingResponseFromJSON(
	jsonBytes []byte,
	bindingResponse *BindingResponse,
) error {
	return json.Unmarshal(jsonBytes, bindingResponse)
}

// ToJSON returns a []byte containing a JSON representation of the binding
// response
func (b *BindingResponse) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}
