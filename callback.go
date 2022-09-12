package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/jsonpointer"
	"github.com/tidwall/gjson"
)

// CallbacksMap is a map of reusable Callback Objects.
type CallbacksMap = ComponentMap[*Callbacks]

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
	Extensions `json:"-"`
	Location   `json:"-"`
	Items      PathItemObjs `json:"-"`
}

// Resolve performs a l
func (c *Callbacks) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return c.resolve(ptr)
}

func (c *Callbacks) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return c, nil
	}
	nxt, tok, _ := ptr.Next()
	item := c.Items.Get(Text(tok))
	if item == nil {
		return nil, newErrNotFound(c.Location.AbsoluteLocation(), tok)
	}
	return item.resolve(nxt)
}

func (c *Callbacks) location() Location {
	return c.Location
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
			c.Items.Set(Text(key.String()), &v)
		}
		return err == nil
	})
	return err
}

// kind returns KindCallback
func (*Callbacks) Kind() Kind     { return KindCallbacks }
func (*Callbacks) mapKind() Kind  { return KindCallbacksMap }
func (Callbacks) sliceKind() Kind { return KindUndefined }

func (c *Callbacks) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = loc
	return c.Items.setLocation(loc)
}

var _ node = (*Callbacks)(nil)
