package openapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

var ErrCallbackNotFound = fmt.Errorf("callback not found")

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
	Location   *Location `json:"-"`
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

// kind returns KindCallback
func (*Callbacks) Kind() Kind { return KindCallbacks }

func (c *Callbacks) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = &loc
	for _, kv := range c.Items {
		if err := kv.Component.Object.setLocation(loc.Append(kv.Key)); err != nil {
			return err
		}
	}
	return nil
}

// CallbackMap is a map of reusable Callback Objects.
type CallbackMap = ComponentMap[*Callbacks]

var _ node = (*Callbacks)(nil)

// func (c *Callbacks) resolve(p string) (node, error) {
// 	ptr, err := jsonpointer.Parse(p)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to resolve p: %q: %w", p, err)
// 	}
// 	np, t, ok := ptr.Next()
// 	if !ok {
// 		return c, nil
// 	}
// 	pi, ok := c.Items.Get(string(t))
// 	if !ok {
// 		return nil, fmt.Errorf("%w: %q", ErrCallbackNotFound, t)
// 	}
// 	return pi.resolve(np.String())
// }
