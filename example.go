package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// ExampleKind indicates wheter the ExampleObj is an Example or a Reference
type ExampleKind uint8

const (
	// ExampleKindObj indicates an ExampleObj
	ExampleKindObj ExampleKind = iota
	// ExampleKindRef indicates a Reference
	ExampleKindRef
)

// Example is either an Example or a Reference
type Example interface {
	ResolveExample(ExampleResolverFunc) (*ExampleObj, error)
	ExampleKind() ExampleKind
}

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
type example ExampleObj

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

// ExampleKind returns ExampleKindObj
func (e *ExampleObj) ExampleKind() ExampleKind { return ExampleKindObj }

// ResolveExample resolves ExampleObj by returning itself. resolve is  not called.
func (e *ExampleObj) ResolveExample(ExampleResolverFunc) (*ExampleObj, error) {
	return e, nil
}

// Examples is an object to hold reusable Examples.
type Examples map[string]Example

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
