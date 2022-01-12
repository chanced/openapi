package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
)

// Callbacks is a map of reusable Callback Objects.
type Callbacks map[string]Callback

// Kind returns KindCallbacks
func (Callbacks) Kind() Kind {
	return KindCallbacks
}

func (cs *Callbacks) Len() int {
	if cs == nil || *cs == nil {
		return 0
	}
	return len(*cs)
}

func (cs *Callbacks) Get(key string) (Callback, bool) {
	if cs == nil || *cs == nil {
		return nil, false
	}
	v, ok := (*cs)[key]
	return v, ok
}

func (cs *Callbacks) Set(key string, val Callback) {
	if *cs == nil {
		*cs = Callbacks{
			key: val,
		}
		return
	}
	(*cs)[key] = val
}

func (cs Callbacks) Nodes() Nodes {
	if cs.Len() == 0 {
		return nil
	}
	m := make(Nodes, cs.Len())
	for k, v := range cs {
		m[k] = NodeDetail{
			Node:       v,
			TargetKind: KindCallback,
		}
	}
	return m
}

// UnmarshalJSON unmarshals JSON
func (c *Callbacks) UnmarshalJSON(data []byte) error {
	var o map[string]json.RawMessage
	res := make(Callbacks, len(o))
	err := json.Unmarshal(data, &o)
	if err != nil {
		return err
	}
	for k, d := range o {
		if isRefJSON(d) {
			v, err := unmarshalReferenceJSON(d)
			if err != nil {
				return err
			}
			res[k] = v
		} else {
			var v CallbackObj
			if err := json.Unmarshal(d, &v); err != nil {
				return err
			}
			res[k] = &v
		}
	}
	*c = res
	return nil
}

// MarshalYAML marshals YAML
func (c Callbacks) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(c)
}

// UnmarshalYAML unmarshals YAML
func (c *Callbacks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, c)
}

type callback CallbackObj

// CallbackObj is map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the path item object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
//
// To describe incoming requests from the API provider independent from another
// API call, use the webhooks field.
type CallbackObj struct {
	Paths      PathItems `json:"-"`
	Extensions `json:"-"`
}

// Nodes returns
func (c *CallbackObj) Nodes() Nodes {
	if c.Paths.Len() == 0 {
		return nil
	}
	m := make(Nodes, len(c.Paths))
	for k, v := range c.Paths {
		m[k] = NodeDetail{
			Node: v,
		}
	}
	return m
}

// MarshalJSON marshals JSON
func (c CallbackObj) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(c.Paths)
	if err != nil {
		return b, err
	}
	return marshalExtendedJSONInto(b, callback(c))
}

// UnmarshalJSON unmarshals JSON
func (c *CallbackObj) UnmarshalJSON(data []byte) error {
	*c = CallbackObj{
		Paths:      PathItems{},
		Extensions: Extensions{},
	}
	var err error
	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		d := []byte(value.Raw)
		if IsExtensionKey(key.String()) {
			c.Extensions.SetEncodedExtension(key.String(), d)
		} else {
			var v Path
			v, err = unmarshalPathJSON(d)
			c.Paths[key.String()] = v
		}
		if err != nil {
			return false
		}
		return true
	})
	return err
}

// MarshalYAML marshals YAML
func (c CallbackObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(c)
}

// UnmarshalYAML unmarshals YAML
func (c *CallbackObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, c)
}

// Kind returns KindCallback
func (*CallbackObj) Kind() Kind {
	return KindCallback
}

// ResolveCallback resolves CallbackObj by returning itself. resolve is  not called.
func (c *CallbackObj) ResolveCallback(func(ref string) (*CallbackObj, error)) (*CallbackObj, error) {
	return c, nil
}

// Callback can either be a CallbackObj or a Reference
type Callback interface {
	Node
	ResolveCallback(func(ref string) (*CallbackObj, error)) (*CallbackObj, error)
}

// ResolvedCallback is map of possible out-of band callbacks related to the parent
// operation. Each value in the map is a Path Item Object that describes a set
// of requests that may be initiated by the API provider and the expected
// responses. The key value used to identify the path item object is an
// expression, evaluated at runtime, that identifies a URL to use for the
// callback operation.
//
// To describe incoming requests from the API provider independent from another
// API call, use the webhooks field.
type ResolvedCallback struct {
	Paths      ResolvedPathItems `json:"-"`
	Extensions `json:"-"`
}

func (rc *ResolvedCallback) Nodes() Nodes {
	if rc == nil {
		return nil
	}
	if rc.Paths.Len() == 0 {
		return nil
	}
	nodes := make(Nodes, len(rc.Paths))
	for k, v := range rc.Paths {
		nodes[k] = NodeDetail{
			Node: v,
		}
	}
	return nodes
}

func (*ResolvedCallback) Kind() Kind {
	return KindResolvedCallback
}

// ResolvedCallbacks is a map of resolved Callback Objects.
type ResolvedCallbacks map[string]*ResolvedCallback

func (ResolvedCallbacks) Kind() Kind {
	return KindResolvedCallbacks
}

var (
	_ Node = (*CallbackObj)(nil)
	_ Node = (*ResolvedCallback)(nil)
	_ Node = (Callbacks)(nil)
	_ Node = (ResolvedCallbacks)(nil)
)
