package openapi

import (
	"encoding/json"

	"github.com/tidwall/gjson"
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
	Paths      PathItems `json:"-"`
	Extensions `json:"-"`
}

type callback Callback

// MarshalJSON marshals JSON
func (c Callback) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(c.Paths)
	if err != nil {
		return b, err
	}
	return marshalExtendedJSONInto(b, callback(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Callback) UnmarshalJSON(data []byte) error {
	*c = Callback{
		Paths:      PathItems{},
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		d := []byte(value.Raw)
		if IsExtensionKey(key.String()) {
			c.Extensions.SetEncodedExtension(key.String(), d)
		} else {
			var comp Component[*Path]
			err = comp.UnmarshalJSON(d)
			if err != nil {
				return false
			}
			c.Paths[key.String()] = comp
		}
		if err != nil {
			return false
		}
		return true
	})
	return err
}

// Kind returns KindCallback
func (Callback) Kind() Kind { return KindCallback }

// Callbacks is a map of reusable Callback Objects.
type Callbacks Map[*Callback]
