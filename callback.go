package openapi

import (
	"encoding/json"
)

// Callback is map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the path item object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
//
// To describe incoming requests from the API provider independent from another
// API call, use the webhooks field.
type Callback struct {
	Paths      PathItemMap `json:"-"`
	Extensions `json:"-"`
}

// MarshalJSON marshals JSON
func (c Callback) MarshalJSON() ([]byte, error) {
	type callback Callback
	b, err := json.Marshal(c.Paths)
	if err != nil {
		return b, err
	}
	return marshalExtendedJSONInto(b, callback(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Callback) UnmarshalJSON(data []byte) error {
	type callback Callback
	var n callback
	err := unmarshalExtendedJSON(data, &n)
	if err != nil {
		return err
	}
	*c = Callback(n)
	return nil
}

// Kind returns KindCallback
func (Callback) Kind() Kind { return KindCallback }

// CallbackMap is a map of reusable Callback Objects.
type CallbackMap = ComponentMap[*Callback]
