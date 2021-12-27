package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
)

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

type callback CallbackObj

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

// Kind returns CallbackKindCallback
func (*CallbackObj) Kind() Kind { return KindCallback }

// ResolveCallback resolves CallbackObj by returning itself. resolve is  not called.
func (c *CallbackObj) ResolveCallback(func(ref string) (*CallbackObj, error)) (*CallbackObj, error) {
	return c, nil
}

// Callback can either be a CallbackObj or a Reference
type Callback interface {
	Node
	ResolveCallback(func(ref string) (*CallbackObj, error)) (*CallbackObj, error)
}

// Callbacks is a map of reusable Callback Objects.
type Callbacks map[string]Callback

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

// ResolvedCallbacks is a map of resolved Callback Objects.
type ResolvedCallbacks map[string]*ResolvedCallback
