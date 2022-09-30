package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonx"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

// ExampleMap is an object to hold reusable ExampleMap.
type ExampleMap = ComponentMap[*Example]

// Example is an example for various api interactions such as Responses
//
// In all cases, the example value is expected to be compatible with the type
// schema of its associated value. Tooling implementations MAY choose to
// validate compatibility automatically, and reject the example value(s) if
// incompatible.
type Example struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// Short description for the example.
	Summary Text `json:"summary,omitempty"`

	// Long description for the example. CommonMark syntax MAY be used for rich
	// text representation.
	Description Text `json:"description,omitempty"`

	// Any embedded literal example. The value field and externalValue field are
	// mutually exclusive. To represent examples of media types that cannot
	// naturally represented in JSON or YAML, use a string value to contain the
	// example, escaping where necessary.
	Value jsonx.RawMessage `json:"value,omitempty"`

	// A URI that points to the literal example. This provides the capability to
	// reference examples that cannot easily be included in JSON or YAML
	// documents. The value field and externalValue field are mutually
	// exclusive. See the rules for resolving Relative References.
	ExternalValue *uri.URI `json:"externalValue,omitempty"`
}

func (e *Example) Nodes() []Node {
	if e == nil {
		return nil
	}
	return downcastNodes(e.nodes())
}
func (e *Example) nodes() []node { return nil }

func (*Example) Refs() []Ref { return nil }

func (e *Example) Anchors() (*Anchors, error) { return nil, nil }

func (*Example) Kind() Kind      { return KindExample }
func (*Example) mapKind() Kind   { return KindExampleMap }
func (*Example) sliceKind() Kind { return KindUndefined }

// MarshalJSON marshals JSON
func (e Example) MarshalJSON() ([]byte, error) {
	type example Example

	return marshalExtendedJSON(example(e))
}

// UnmarshalJSON unmarshals JSON
func (e *Example) UnmarshalJSON(data []byte) error {
	type example Example
	var v example
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*e = Example(v)
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (e Example) MarshalYAML() (interface{}, error) {
	j, err := e.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(j, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (e *Example) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, e)
}

func (e *Example) setLocation(loc Location) error {
	if e == nil {
		return nil
	}
	e.Location = loc
	return nil
}
func (e *Example) isNil() bool { return e == nil }

func (*Example) refable() {}

var (
	_ node = (*Example)(nil)

	_ node = (*ExampleMap)(nil)
)

// func (e *Example) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return e.resolveNodeByPointer(ptr)
// }

// func (e *Example) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return e, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(e.Location.AbsoluteLocation(), tok)
// }
