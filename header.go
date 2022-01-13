package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
)

// Header is either a Header or a Reference
type Header interface {
	Node
	ResolveHeader(func(ref string) (*HeaderObj, error)) (*HeaderObj, error)
}

// HeaderObj follows the structure of the Parameter Object with the following
// changes:
//		- name MUST NOT be specified, it is given in the corresponding headers map.
//		- in MUST NOT be specified, it is implicitly in header.
//		- All traits that are affected by the location MUST be applicable to a
// 		  location of header (for example, style).
type HeaderObj struct {
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
	Schema *SchemaObj `json:"schema,omitempty"`
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

func (h *HeaderObj) Nodes() Nodes {
	n := makeNodes(nodes{
		{"schema", h.Schema, KindSchema},
		{"examples", h.Examples, KindExamples},
	})
	if len(n) == 0 {
		return nil
	}
	return n
}

type header HeaderObj

// Kind returns KindHeader
func (*HeaderObj) Kind() Kind {
	return KindHeader
}

// ResolveHeader resolves HeaderObj by returning itself. resolve is  not called.
func (h *HeaderObj) ResolveHeader(func(ref string) (*HeaderObj, error)) (*HeaderObj, error) {
	return h, nil
}

// MarshalJSON marshals h into JSON
func (h HeaderObj) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(header(h))
}

// UnmarshalJSON unmarshals json into h
func (h *HeaderObj) UnmarshalJSON(data []byte) error {
	v := header{}
	err := unmarshalExtendedJSON(data, &v)
	*h = HeaderObj(v)
	return err
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (h HeaderObj) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// UnmarshalYAML unmarshals yaml into s
func (h *HeaderObj) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, h)
}

// Headers holds reusable HeaderObjs.
type Headers map[string]Header

func (es *Headers) Len() int {
	if es == nil || *es == nil {
		return 0
	}
	return len(*es)
}

func (es *Headers) Get(key string) (Header, bool) {
	if es == nil || *es == nil {
		return nil, false
	}
	v, ok := (*es)[key]
	return v, ok
}

func (es *Headers) Set(key string, val Header) {
	if *es == nil {
		*es = Headers{
			key: val,
		}
		return
	}
	(*es)[key] = val
}

func (es Headers) Nodes() Nodes {
	if len(es) == 0 {
		return nil
	}
	nodes := make(Nodes, len(es))
	for k, v := range es {
		nodes[k] = NodeDetail{
			TargetKind: KindHeader,
			Node:       v,
		}
	}
	return nodes
}

// Kind returns KindHeaders
func (Headers) Kind() Kind {
	return KindHeaders
}

// UnmarshalJSON unmarshals JSON data into p
func (h *Headers) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	*h = make(Headers, len(m))
	if err != nil {
		return err
	}
	for i, j := range m {
		if isRefJSON(data) {
			var v Reference
			if err = json.Unmarshal(j, &v); err != nil {
				return err
			}
			(*h)[i] = &v
		} else {
			var v HeaderObj
			if err = json.Unmarshal(j, &v); err != nil {
				return err
			}
			(*h)[i] = &v
		}
	}
	return nil
}

// UnmarshalYAML unmarshals YAML data into p
func (h *Headers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, h)
}

// MarshalYAML marshals p into YAML
func (h Headers) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// ResolvedHeaders holds reusable ResolvedHeaders.
type ResolvedHeaders map[string]*ResolvedHeader

// Kind returns KindResolvedHeaders
func (ResolvedHeaders) Kind() Kind {
	return KindResolvedHeaders
}

// ResolvedHeader follows the structure of the Parameter Object with the following
// changes:
//		- name MUST NOT be specified, it is given in the corresponding headers map.
//		- in MUST NOT be specified, it is implicitly in header.
//		- All traits that are affected by the location MUST be applicable to a
// 		  location of header (for example, style).
//
// ResolvedHeader has been fully resolved.
type ResolvedHeader struct {

	// TODO: add reference

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
	Schema *ResolvedSchema `json:"schema,omitempty"`
	// Examples of the parameter's potential value. Each example SHOULD
	// contain a value in the correct format as specified in the parameter
	// encoding. The examples field is mutually exclusive of the example
	// field. Furthermore, if referencing a schema that contains an example,
	// the examples value SHALL override the example provided by the schema.
	Examples ResolvedExamples `json:"examples,omitempty"`
	// OpenAPI extensions
	Extensions `json:"-"`
}

func (rh *ResolvedHeader) Nodes() Nodes {
	n := makeNodes(nodes{
		{"schema", rh.Schema, KindResolvedSchema},
		{"examples", rh.Examples, KindResolvedExamples},
	})
	if len(n) == 0 {
		return nil
	}
	return n
}

// Kind returns KindResolvedHeader
func (*ResolvedHeader) Kind() Kind {
	return KindResolvedHeader
}

func (rhs *ResolvedHeaders) Len() int {
	if rhs == nil || *rhs == nil {
		return 0
	}
	return len(*rhs)
}

func (rhs *ResolvedHeaders) Get(key string) (*ResolvedHeader, bool) {
	if rhs == nil || *rhs == nil {
		return nil, false
	}
	v, ok := (*rhs)[key]
	return v, ok
}

func (rhs *ResolvedHeaders) Set(key string, val *ResolvedHeader) {
	if *rhs == nil {
		*rhs = ResolvedHeaders{
			key: val,
		}
		return
	}
	(*rhs)[key] = val
}

func (rhs ResolvedHeaders) Nodes() Nodes {
	if len(rhs) == 0 {
		return nil
	}
	nodes := make(Nodes, len(rhs))
	for k, v := range rhs {
		nodes[k] = NodeDetail{
			TargetKind: KindResolvedHeader,
			Node:       v,
		}
	}
	return nodes
}

var (
	_ Node = (*HeaderObj)(nil)
	_ Node = (Headers)(nil)
	_ Node = (*ResolvedHeader)(nil)
	_ Node = (ResolvedHeaders)(nil)
)
