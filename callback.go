package openapi

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

type CallbackItemMap = ComponentMap[*PathItem]

// Callbacks is map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the path item object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
//
// To describe incoming requests from the API provider independent from another
// API call, use the webhooks field.
type Callbacks struct {
	Items      PathItemMap `json:"-"`
	Extensions `json:"-"`
}

// MarshalJSON marshals JSON
func (c Callbacks) MarshalJSON() ([]byte, error) {
	type callback Callbacks
	b, err := json.Marshal(c.Items)
	if err != nil {
		return b, err
	}
	return marshalExtendedJSONInto(b, callback(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Callbacks) UnmarshalJSON(data []byte) error {
	*c = Callbacks{
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			c.SetRawExtension(key.String(), []byte(value.Raw))
		} else {
			var v PathItem
			err = json.Unmarshal([]byte(value.Raw), &v)
			c.Items = append(c.Items, ComponentEntry[*PathItem]{
				Key:       key.String(),
				Component: Component[*PathItem]{Object: &v},
			})
		}
		return err == nil
	})
	return err
}

// kind returns kindCallback
func (*Callbacks) kind() kind { return kindCallbacks }

// CallbackMap is a map of reusable Callback Objects.
type CallbackMap = ComponentMap[*Callbacks]
