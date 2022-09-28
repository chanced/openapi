package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// EncodingMap is a ComponentMap between a property name and its encoding information. The
// key, being the property name, MUST exist in the schema as a property. The
// encoding object SHALL only apply to requestBody objects when the media type
// is multipart or application/x-www-form-urlencoded.
type EncodingMap = ComponentMap[*Encoding]

// Encoding definition applied to a single schema property.
type Encoding struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The Content-Type for encoding a specific property. Default value depends
	// on the property type:
	//
	//  - for object - application/json;
	//  - for array – the default is defined based on the inner type;
	//  - for all other cases the default is application/octet-stream.
	// The value can be a specific media type (e.g. application/json), a
	// wildcard media type (e.g. image/*), or a comma-separated list of the two
	// types.
	ContentType Text `json:"contentType,omitempty"`
	// A map allowing additional information to be provided as headers, for
	// example Content-Disposition. Content-Type is described separately and
	// SHALL be ignored in this section. This property SHALL be ignored if the
	// request body media type is not a multipart.
	Headers *HeaderMap `json:"headers,omitempty"`
	// Describes how a specific property value will be serialized depending on
	// its type. See Parameter Object for details on the style property. The
	// behavior follows the same values as query parameters, including default
	// values. This property SHALL be ignored if the request body media type is
	// not application/x-www-form-urlencoded or multipart/form-data. If a value
	// is explicitly defined, then the value of contentType (implicit or
	// explicit) SHALL be ignored.
	Style Text `json:"style,omitempty"`
	// When this is true, property values of type array or object generate
	// separate parameters for each value of the array, or key-value-pair of the
	// map. For other types of properties this property has no effect. When
	// style is form, the default value is true. For all other styles, the
	// default value is false. This property SHALL be ignored if the request
	// body media type is not application/x-www-form-urlencoded or
	// multipart/form-data. If a value is explicitly defined, then the value of
	// contentType (implicit or explicit) SHALL be ignored.
	Explode *bool `json:"explode,omitempty"`
	// Determines whether the parameter value SHOULD allow reserved characters,
	// as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without
	// percent-encoding. The default value is false. This property SHALL be
	// ignored if the request body media type is not
	// application/x-www-form-urlencoded or multipart/form-data. If a value is
	// explicitly defined, then the value of contentType (implicit or explicit)
	// SHALL be ignored.
	AllowReserved *bool `json:"allowReserved,omitempty"`
}

func (e *Encoding) Nodes() []Node {
	if e == nil {
		return nil
	}
	return downcastNodes(e.nodes())
}

func (e *Encoding) nodes() []node {
	if e == nil {
		return nil
	}
	return appendEdges(nil, e.Headers)
}

func (e *Encoding) Refs() []Ref {
	if e == nil {
		return nil
	}
	return e.Headers.Refs()
}

func (e *Encoding) Anchors() (*Anchors, error) {
	if e == nil {
		return nil, nil
	}
	return e.Headers.Anchors()
}

// func (e *Encoding) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	err := ptr.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return e.resolveNodeByPointer(ptr)
// }

// func (e *Encoding) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return e, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch nxt {
// 	case "headers":
// 		if e.Headers == nil {
// 			return nil, newErrNotFound(e.Location.AbsoluteLocation(), tok)
// 		}
// 		return e.Headers.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(e.Location.AbsoluteLocation(), tok)
// 	}
// }

func (*Encoding) Kind() Kind      { return KindEncoding }
func (*Encoding) mapKind() Kind   { return KindEncodingMap }
func (*Encoding) sliceKind() Kind { return KindUndefined }

func (e *Encoding) setLocation(loc Location) error {
	if e == nil {
		return nil
	}
	e.Location = loc
	return e.Headers.setLocation(loc.AppendLocation("headers"))
}

// MarshalJSON marshals e into JSON
func (e Encoding) MarshalJSON() ([]byte, error) {
	type encoding Encoding
	return marshalExtendedJSON(encoding(e))
}

// UnmarshalJSON unmarshals json into e
func (e *Encoding) UnmarshalJSON(data []byte) error {
	type encoding Encoding
	v := encoding{}
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*e = Encoding(v)
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (e Encoding) MarshalYAML() (interface{}, error) {
	j, err := e.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (e *Encoding) UnmarshalYAML(value *yaml.Node) error {
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

func (e *Encoding) isNil() bool { return e == nil }

func (*Encoding) refable() {}

var _ node = (*Encoding)(nil)
