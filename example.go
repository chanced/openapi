package openapi

import (
	"encoding/json"
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

// kind returns kindExample
func (*Example) kind() kind { return kindExample }
