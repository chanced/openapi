package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
)

// HeaderMap holds reusable HeaderMap.
type HeaderMap = ComponentMap[*Header]

// Header follows the structure of the Parameter Object with the following
// changes:
//   - name MUST NOT be specified, it is given in the corresponding headers map.
//   - in MUST NOT be specified, it is implicitly in header.
//   - All traits that are affected by the location MUST be applicable to a
//     location of header (for example, style).
type Header struct {
	// OpenAPI extensions
	Extensions `json:"-"`
	Location   `json:"-"`

	// A brief description of the parameter. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description Text `json:"description,omitempty"`
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
	Style Style `json:"style,omitempty"`
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
	Examples *ExampleMap `json:"examples,omitempty"`
	// Example of the parameter's potential value. The example SHOULD match the
	// specified schema and encoding properties if present. The example field is
	// mutually exclusive of the examples field. Furthermore, if referencing a
	// schema that contains an example, the example value SHALL override the
	// example provided by the schema. To represent examples of media types that
	// cannot naturally be represented in JSON or YAML, a string value can
	// contain the example with escaping where necessary.
	Example jsonx.RawMessage `json:"example,omitempty"`
}

func (h *Header) Anchors() (*Anchors, error) {
	if h == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = h.Schema.Anchors(); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(h.Examples.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

func (h *Header) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return h.resolveNodeByPointer(ptr)
}

func (h *Header) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return h, nil
	}
	nxt, tok, _ := ptr.Next()
	switch nxt {
	case "schema":
		if h.Schema == nil {
			return nil, newErrNotFound(h.Location.AbsoluteLocation(), tok)
		}
		return h.Schema.resolveNodeByPointer(nxt)
	case "examples":
		if h.Examples == nil {
			return nil, newErrNotFound(h.Location.AbsoluteLocation(), tok)
		}
		return h.Examples.resolveNodeByPointer(nxt)
	default:
		return nil, newErrNotResolvable(h.Location.AbsoluteLocation(), tok)
	}
}

func (*Header) Kind() Kind      { return KindHeader }
func (*Header) mapKind() Kind   { return KindHeaderMap }
func (*Header) sliceKind() Kind { return KindUndefined }

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

func (h *Header) setLocation(loc Location) error {
	h.Location = loc
	if err := h.Examples.setLocation(loc.Append("examples")); err != nil {
		return err
	}

	if err := h.Schema.setLocation(loc.Append("schema")); err != nil {
		return err
	}
	return nil
}
func (h *Header) isNil() bool { return h == nil }

var _ node = (*Header)(nil)
