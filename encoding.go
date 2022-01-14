package openapi

import (
	"github.com/chanced/openapi/yamlutil"
)

// Encodings is a map between a property name and its encoding information. The
// key, being the property name, MUST exist in the schema as a property. The
// encoding object SHALL only apply to requestBody objects when the media type
// is multipart or application/x-www-form-urlencoded.
type Encodings map[string]*Encoding

// Kind returns KindEncodings
func (Encodings) Kind() Kind {
	return KindEncodings
}

func (es *Encodings) Len() int {
	if es == nil || *es == nil {
		return 0
	}
	return len(*es)
}

func (es *Encodings) Get(key string) (*Encoding, bool) {
	if es == nil || *es == nil {
		return nil, false
	}
	v, ok := (*es)[key]
	return v, ok
}

func (es *Encodings) Set(key string, val *Encoding) {
	if *es == nil {
		*es = Encodings{
			key: val,
		}
		return
	}
	(*es)[key] = val
}

func (es Encodings) Nodes() Nodes {
	if len(es) == 0 {
		return nil
	}
	n := make(Nodes, len(es))
	for k, v := range es {
		n.maybeAdd(k, v, KindEncoding)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

// Encoding definition applied to a single schema property.
type Encoding struct {
	// The Content-Type for encoding a specific property. Default value depends
	// on the property type:
	//
	//  - for object - application/json;
	//  - for array – the default is defined based on the inner type;
	//  - for all other cases the default is application/octet-stream.
	// The value can be a specific media type (e.g. application/json), a
	// wildcard media type (e.g. image/*), or a comma-separated list of the two
	// types.
	ContentType string `json:"contentType,omitempty"`
	// A map allowing additional information to be provided as headers, for
	// example Content-Disposition. Content-Type is described separately and
	// SHALL be ignored in this section. This property SHALL be ignored if the
	// request body media type is not a multipart.
	Headers Headers `json:"headers,omitempty"`
	// Describes how a specific property value will be serialized depending on
	// its type. See Parameter Object for details on the style property. The
	// behavior follows the same values as query parameters, including default
	// values. This property SHALL be ignored if the request body media type is
	// not application/x-www-form-urlencoded or multipart/form-data. If a value
	// is explicitly defined, then the value of contentType (implicit or
	// explicit) SHALL be ignored.
	Style Style `json:"style,omitempty"`
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

	Extensions `json:"-"`
}
type encodingobj Encoding

func (e *Encoding) Nodes() Nodes {
	return makeNodes(nodes{
		{"headers", e.Headers, KindHeaders},
	})
}

// Kind returns KindEncoding
func (*Encoding) Kind() Kind {
	return KindEncoding
}

// MarshalJSON marshals e into JSON
func (e Encoding) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(encodingobj(e))
}

// UnmarshalJSON unmarshals json into e
func (e *Encoding) UnmarshalJSON(data []byte) error {
	v := encodingobj{}
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*e = Encoding(v)
	return nil
}

// MarshalYAML marshals YAML
func (e Encoding) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(e)
}

// UnmarshalYAML unmarshals YAML
func (e *Encoding) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, e)
}

// ResolvedEncodings is a map between a property name and its encoding
// information. The key, being the property name, MUST exist in the schema as a
// property. The encoding object SHALL only apply to requestBody objects when
// the media type is multipart or application/x-www-form-urlencoded.
type ResolvedEncodings map[string]*ResolvedEncoding

func (res *ResolvedEncodings) Len() int {
	if res == nil || *res == nil {
		return 0
	}
	return len(*res)
}

func (res *ResolvedEncodings) Get(key string) (*ResolvedEncoding, bool) {
	if res == nil || *res == nil {
		return nil, false
	}
	v, ok := (*res)[key]
	return v, ok
}

func (res *ResolvedEncodings) Set(key string, val *ResolvedEncoding) {
	if *res == nil {
		*res = ResolvedEncodings{
			key: val,
		}
		return
	}
	(*res)[key] = val
}

func (res ResolvedEncodings) Nodes() Nodes {
	if len(res) == 0 {
		return nil
	}
	n := make(Nodes, len(res))
	for k, v := range res {
		n.maybeAdd(k, v, KindResolvedEncoding)
	}
	if len(n) == 0 {
		return nil
	}
	return n
}

// Kind returns KindResolvedEncodings
func (ResolvedEncodings) Kind() Kind {
	return KindResolvedEncodings
}

// ResolvedEncoding definition applied to a single schema property.
type ResolvedEncoding struct {
	// The Content-Type for encoding a specific property. Default value depends
	// on the property type:
	//
	//  - for object - application/json;
	//  - for array – the default is defined based on the inner type;
	//  - for all other cases the default is application/octet-stream.
	// The value can be a specific media type (e.g. application/json), a
	// wildcard media type (e.g. image/*), or a comma-separated list of the two
	// types.
	ContentType string `json:"contentType,omitempty"`
	// A map allowing additional information to be provided as headers, for
	// example Content-Disposition. Content-Type is described separately and
	// SHALL be ignored in this section. This property SHALL be ignored if the
	// request body media type is not a multipart.
	Headers ResolvedHeaders `json:"headers,omitempty"`
	// Describes how a specific property value will be serialized depending on
	// its type. See Parameter Object for details on the style property. The
	// behavior follows the same values as query parameters, including default
	// values. This property SHALL be ignored if the request body media type is
	// not application/x-www-form-urlencoded or multipart/form-data. If a value
	// is explicitly defined, then the value of contentType (implicit or
	// explicit) SHALL be ignored.
	Style Style `json:"style,omitempty"`
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

	Extensions `json:"-"`
}

// Kind returns KindResolvedEncoding
func (*ResolvedEncoding) Kind() Kind {
	return KindResolvedEncoding
}

func (e *ResolvedEncoding) Nodes() Nodes {
	return makeNodes(nodes{
		{"headers", e.Headers, KindResolvedHeaders},
	})
}

var (
	_ Node = (*Encoding)(nil)
	_ Node = (Encodings)(nil)
	_ Node = (*ResolvedEncoding)(nil)
	_ Node = (ResolvedEncodings)(nil)
)
