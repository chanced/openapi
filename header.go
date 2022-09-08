package openapi

import (
	"encoding/json"
)

// Headers holds reusable Headers.
type Headers Map[*Header]

// Header follows the structure of the Parameter Object with the following
// changes:
//   - name MUST NOT be specified, it is given in the corresponding headers map.
//   - in MUST NOT be specified, it is implicitly in header.
//   - All traits that are affected by the location MUST be applicable to a
//     location of header (for example, style).
type Header struct {
	// A brief description of the parameter. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty"`
	// Determines whether this parameter is mandatory. If the parameter location
	// is "path", this property is REQUIRED and its value MUST be true.
	// Otherwise, the property MAY be included and its default value is false.
	Required *bool `json:"required,omitempty"`
	// Specifies that a parameter is deprecated and SHOULD be transitioned out
	// of usage. Default value is false.
	Deprecated *bool `json:"deprecated,omitempty"`
	// Sets the ability to pass empty-valued parameters. This is valid only for
	// query parameters and allows sending a parameter with an empty value.
	// Default value is false. If style is used, and if behavior is n/a (cannot
	// be serialized), the value of allowEmptyValue SHALL be ignored. Use of
	// this property is NOT RECOMMENDED, as it is likely to be removed in a
	// later revision.
	AllowEmptyValue *bool `json:"allowEmptyValue,omitempty"`
	// Describes how the parameter value will be serialized depending on the
	// type of the parameter value.
	// Default values (based on value of in):
	// 	- for query - form;
	// 	- for path - simple;
	// 	- for header - simple;
	// 	- for cookie - form.
	Style string `json:"style,omitempty"`
	// When this is true, parameter values of type array or object generate
	// separate parameters for each value of the array or key-value pair of the
	// map. For other types of parameters this property has no effect. When
	// style is form, the default value is true. For all other styles, the
	// default value is false.
	Explode *bool `json:"explode,omitempty"`
	// Determines whether the parameter value SHOULD allow reserved characters,
	// as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without
	// percent-encoding. This property only applies to parameters with an in
	// value of query. The default value is false.
	AllowReserved *bool `json:"allowReserved,omitempty"`
	// The schema defining the type used for the parameter.
	Schema *Schema `json:"schema,omitempty"`
	// Examples of the parameter's potential value. Each example SHOULD
	// contain a value in the correct format as specified in the parameter
	// encoding. The examples field is mutually exclusive of the example
	// field. Furthermore, if referencing a schema that contains an example,
	// the examples value SHALL override the example provided by the schema.
	Examples Examples `json:"examples,omitempty"`
	// Example of the parameter's potential value. The example SHOULD match the
	// specified schema and encoding properties if present. The example field is
	// mutually exclusive of the examples field. Furthermore, if referencing a
	// schema that contains an example, the example value SHALL override the
	// example provided by the schema. To represent examples of media types that
	// cannot naturally be represented in JSON or YAML, a string value can
	// contain the example with escaping where necessary.
	Example json.RawMessage `json:"example,omitempty"`
	// OpenAPI extensions
	Extensions `json:"-"`
}

// MarshalJSON marshals h into JSON
func (h Header) MarshalJSON() ([]byte, error) {
	type header Header

	return marshalExtendedJSON(header(h))
}

// UnmarshalJSON unmarshals json into h
func (h *Header) UnmarshalJSON(data []byte) error {
	type header Header

	v := header{}
	err := unmarshalExtendedJSON(data, &v)
	*h = Header(v)
	return err
}
func (Header) Kind() Kind { return KindHeader }
