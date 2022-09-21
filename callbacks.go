package openapi

import (
	"encoding/json"
	"strings"

	"github.com/chanced/transcode"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
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
	PathItems  `json:"-"`
}

func (c *Callbacks) Nodes() []Node {
	if c == nil {
		return nil
	}
	return downcastNodes(c.nodes())
}

func (c *Callbacks) nodes() []node {
	if c == nil {
		return nil
	}
	edges := make([]node, 0, 1)
	edges = appendEdges(edges, c.PathItems.nodes()...)
	return edges
}

func (c *Callbacks) ref() Ref { return nil }

// Edges returns the immediate edges of the Node. This is used to build a
// graph of the OpenAPI document.
//

// kind returns KindCallback
func (*Callbacks) Kind() Kind     { return KindCallbacks }
func (*Callbacks) mapKind() Kind  { return KindCallbacksMap }
func (Callbacks) sliceKind() Kind { return KindUndefined }
func (c *Callbacks) isNil() bool {
	return c == nil
}

func (c *Callbacks) Anchors() (*Anchors, error) {
	if c == nil {
		return nil, nil
	}
	return c.PathItems.Anchors()
}

func (c *Callbacks) Refs() []Ref {
	if c == nil {
		return nil
	}
	return c.PathItems.Refs()
}

// // ResolveNodeByPointer performs a l
// func (c *Callbacks) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return c.resolveNodeByPointer(ptr)
// }

// func (c *Callbacks) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return c, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	item := c.Items.Get(Text(tok))
// 	if item == nil {
// 		return nil, newErrNotFound(c.Location.AbsoluteLocation(), tok)
// 	}
// 	return item.resolveNodeByPointer(nxt)
// }

func (c *Callbacks) location() Location {
	return c.Location
}

// MarshalJSON marshals JSON
func (c Callbacks) MarshalJSON() ([]byte, error) {
	type callback Callbacks
	return marshalExtendedJSON(callback(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Callbacks) UnmarshalJSON(data []byte) error {
	*c = Callbacks{
		Extensions: Extensions{},
	}

	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		if strings.HasPrefix(key.String(), "x-") {
			c.SetRawExtension(Text(key.String()), []byte(value.Raw))
		} else {
			var v PathItem
			err = json.Unmarshal([]byte(value.Raw), &v)
			c.Set(Text(key.String()), &v)
		}
		return err == nil
	})
	return err
}

func (c *Callbacks) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = loc
	return c.PathItems.setLocation(loc)
}

func (c Callbacks) MarshalYAML() (interface{}, error) {
	j, err := c.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Callbacks) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, c)
}

// func (c *Callbacks) Walk(v Visitor) error {
// 	if v == nil {
// 		return nil
// 	}
// 	v, err := v.Visit(c)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}

//		return c.Items.Walk(v)
//	}
func (c *Callbacks) refable() {}

var _ node = (*Callbacks)(nil)
