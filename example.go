package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Example is either an Example or a Reference
type Example interface {
	Node
	ResolveExample(func(ref string) (*ExampleObj, error)) (*ExampleObj, error)
}

// Examples is an object to hold reusable Examples.
type Examples map[string]Example

func (Examples) Kind() Kind {
	return KindExamples
}

type example ExampleObj

// ExampleObj is an example for various api interactions such as Responses
//
// In all cases, the example value is expected to be compatible with the type
// schema of its associated value. Tooling implementations MAY choose to
// validate compatibility automatically, and reject the example value(s) if
// incompatible.
type ExampleObj struct {
	// Short description for the example.
	Summary string `json:"summary,omitempty"`
	// Long description for the example. CommonMark syntax MAY be used for rich
	// text representation.
	Description string `json:"description,omitempty"`
	// Any embedded literal example. The value field and externalValue field are
	// mutually exclusive. To represent examples of media types that cannot
	// naturally represented in JSON or YAML, use a string value to contain the
	// example, escaping where necessary.
	Value json.RawMessage `json:"value,omitempty"`
	// A URI that points to the literal example. This provides the capability to
	// reference examples that cannot easily be included in JSON or YAML
	// documents. The value field and externalValue field are mutually
	// exclusive. See the rules for resolving Relative References.
	ExternalValue string `json:"externalValue,omitempty"`
	Extensions    `json:"-"`
}

// Kind returns KindExample
func (*ExampleObj) Kind() Kind {
	return KindExample
}

// MarshalJSON marshals JSON
func (e ExampleObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(example(e))
}

// UnmarshalJSON unmarshals JSON
func (e *ExampleObj) UnmarshalJSON(data []byte) error {
	var v example
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*e = ExampleObj(v)
	return nil
}

// MarshalYAML marshals YAML
func (e ExampleObj) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(e)
}

// UnmarshalYAML unmarshals YAML
func (e *ExampleObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, e)
}

// Kind returns KindExample
func (e *ExampleObj) ExampleKind() Kind { return KindExample }

// ResolveExample resolves ExampleObj by returning itself. resolve is  not called.
func (e *ExampleObj) ResolveExample(func(ref string) (*ExampleObj, error)) (*ExampleObj, error) {
	return e, nil
}

// UnmarshalJSON unmarshals JSON
func (e *Examples) UnmarshalJSON(data []byte) error {
	var dm map[string]json.RawMessage
	if err := json.Unmarshal(data, &dm); err != nil {
		return err
	}
	res := make(Examples, len(dm))
	for k, d := range dm {
		if isRefJSON(d) {
			v, err := unmarshalReferenceJSON(d)
			if err != nil {
				return err
			}
			res[k] = v
			continue
		}
		var v example
		if err := unmarshalExtendedJSON(d, &v); err != nil {
			return err
		}
		ev := ExampleObj(v)
		res[k] = &ev
	}
	*e = res
	return nil
}

// ResolvedExamples is a map of *ResolvedExamples
type ResolvedExamples map[string]*ResolvedExample

// Kind returns KindResolvedExamples
func (ResolvedExamples) Kind() Kind {
	return KindResolvedExamples
}

// ResolvedExample is an example for various api interactions such as Responses
// that has been resolved.
//
// In all cases, the example value is expected to be compatible with the type
// schema of its associated value. Tooling implementations MAY choose to
// validate compatibility automatically, and reject the example value(s) if
// incompatible.
type ResolvedExample struct {

	// TODO: Add reference

	// Short description for the example.
	Summary string `json:"summary,omitempty"`
	// Long description for the example. CommonMark syntax MAY be used for rich
	// text representation.
	Description string `json:"description,omitempty"`
	// Any embedded literal example. The value field and externalValue field are
	// mutually exclusive. To represent examples of media types that cannot
	// naturally represented in JSON or YAML, use a string value to contain the
	// example, escaping where necessary.
	Value json.RawMessage `json:"value,omitempty"`
	// A URI that points to the literal example. This provides the capability to
	// reference examples that cannot easily be included in JSON or YAML
	// documents. The value field and externalValue field are mutually
	// exclusive. See the rules for resolving Relative References.
	ExternalValue string `json:"externalValue,omitempty"`
	Extensions    `json:"-"`
}

// Kind returns KindResolvedExample
func (*ResolvedExample) Kind() Kind {
	return KindResolvedExample
}

var _ Node = (*ExampleObj)(nil)
var _ Node = (Examples)(nil)
var _ Node = (ResolvedExamples)(nil)
var _ Node = (*ResolvedExample)(nil)
