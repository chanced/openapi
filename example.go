package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
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
	ExternalValue Text `json:"externalValue,omitempty"`

	Extensions `json:"-"`
	Location   `json:"-"`
}

func (e *Example) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return e.resolve(ptr)
}

func (e *Example) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return e, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(e.Location.AbsoluteLocation(), tok)
}

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

func (e *Example) setLocation(loc Location) error {
	if e == nil {
		return nil
	}
	e.Location = loc
	return nil
}

var (
	_ node = (*Example)(nil)
	_ node = (*ExampleMap)(nil)
)
