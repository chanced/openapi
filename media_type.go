package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
)

// ContentMap / MediaTypeMap is a map containing descriptions of potential response payloads. The key is
// a media type or media type range and the value describes it. For
// responses that match multiple keys, only the most specific key is
// applicable. e.g. text/plain overrides text/*
type (
	ContentMap   = ObjMap[*MediaType]
	MediaTypeMap = ObjMap[*MediaType]
)

// MediaType  provides schema and examples for the media type identified by its
// key.
type MediaType struct {
	Extensions `json:"-"`
	Location   `json:"-"`

	// The schema defining the content of the request, response, or parameter.
	Schema *Schema `json:"schema,omitempty"`
	// Example of the media type. The example object SHOULD be in the correct
	// format as specified by the media type. The example field is mutually
	// exclusive of the examples field. Furthermore, if referencing a schema
	// which contains an example, the example value SHALL override the example
	// provided by the schema.
	Example jsonx.RawMessage `json:"example,omitempty"`
	// Examples of the media type. Each example object SHOULD match the media
	// type and specified schema if present. The examples field is mutually
	// exclusive of the example field. Furthermore, if referencing a schema
	// which contains an example, the examples value SHALL override the example
	// provided by the schema.
	Examples *ExampleMap `json:"examples,omitempty"`
	// A map between a property name and its encoding information. The key,
	// being the property name, MUST exist in the schema as a property. The
	// encoding object SHALL only apply to requestBody objects when the media
	// type is multipart or application/x-www-form-urlencoded.
	Encoding *EncodingMap `json:"encoding,omitempty"`
}

func (mt *MediaType) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return mt.resolve(ptr)
}

func (mt *MediaType) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return mt, nil
	}
	nxt, tok, _ := ptr.Next()
	var n node
	switch tok {
	case "schema":
		n = mt.Schema
	case "examples":
		n = mt.Examples
	case "encoding":
		n = mt.Encoding
	default:
		return nil, newErrNotResolvable(mt.Location.AbsoluteLocation(), tok)
	}
	if nxt.IsRoot() {
		return n, nil
	}

	if n == nil {
		return nil, newErrNotFound(mt.Location.AbsoluteLocation(), tok)
	}
	return n.resolve(nxt)
}

// MarshalJSON marshals mt into JSON
func (mt MediaType) MarshalJSON() ([]byte, error) {
	type mediatype MediaType
	return marshalExtendedJSON(mediatype(mt))
}

// UnmarshalJSON unmarshals json into mt
func (mt *MediaType) UnmarshalJSON(data []byte) error {
	type mediatype MediaType

	var v mediatype
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*mt = MediaType(v)
	return nil
}

func (mt *MediaType) setLocation(loc Location) error {
	if mt == nil {
		return nil
	}
	mt.Location = loc
	if err := mt.Schema.setLocation(loc.Append("schema")); err != nil {
		return err
	}
	if err := mt.Examples.setLocation(loc.Append("examples")); err != nil {
		return err
	}
	if err := mt.Encoding.setLocation(loc.Append("encoding")); err != nil {
		return err
	}

	return nil
}
func (*MediaType) Kind() Kind      { return KindMediaType }
func (*MediaType) mapKind() Kind   { return KindMediaTypeMap }
func (*MediaType) sliceKind() Kind { return KindUndefined }

var _ node = (*MediaType)(nil)
