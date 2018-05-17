package service

import "encoding/json"

// Parameters ...
// TODO: krancour: Document this
type Parameters struct {
	Data map[string]interface{}
}

// MarshalJSON ...
// TODO: krancour: Document this
func (p Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Data)
}

// UnmarshalJSON ...
// TODO: krancour: Document this
func (p *Parameters) UnmarshalJSON(bytes []byte) error {
	p.Data = map[string]interface{}{}
	return json.Unmarshal(bytes, &p.Data)
}
