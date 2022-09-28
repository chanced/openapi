package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonx"
	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
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

func (mt *MediaType) Nodes() []Node {
	if mt == nil {
		return nil
	}
	return downcastNodes(mt.nodes())
}

func (mt *MediaType) nodes() []node {
	if mt == nil {
		return nil
	}
	return appendEdges(nil, mt.Schema, mt.Examples, mt.Encoding)
}

func (mt *MediaType) Refs() []Ref {
	if mt == nil {
		return nil
	}
	refs := mt.Schema.Refs()
	refs = append(refs, mt.Examples.Refs()...)
	refs = append(refs, mt.Encoding.Refs()...)
	return refs
}

func (mt *MediaType) Anchors() (*Anchors, error) {
	if mt == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = mt.Schema.Anchors(); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(mt.Examples.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(mt.Encoding.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

// func (mt *MediaType) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return mt.resolveNodeByPointer(ptr)
// }

// func (mt *MediaType) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return mt, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "schema":
// 		if mt.Schema == nil {
// 			return nil, newErrNotFound(mt.AbsoluteLocation(), tok)
// 		}
// 		return mt.Schema.resolveNodeByPointer(nxt)
// 	case "examples":
// 		if mt.Examples == nil {
// 			return nil, newErrNotFound(mt.AbsoluteLocation(), tok)
// 		}
// 		return mt.Examples.resolveNodeByPointer(nxt)
// 	case "encoding":
// 		if mt.Encoding == nil {
// 			return nil, newErrNotFound(mt.AbsoluteLocation(), tok)
// 		}
// 		return mt.Encoding.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(mt.Location.AbsoluteLocation(), tok)

// 	}
// }

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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (mt MediaType) MarshalYAML() (interface{}, error) {
	j, err := mt.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (mt *MediaType) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, mt)
}

func (mt *MediaType) setLocation(loc Location) error {
	if mt == nil {
		return nil
	}
	mt.Location = loc
	if err := mt.Schema.setLocation(loc.AppendLocation("schema")); err != nil {
		return err
	}
	if err := mt.Examples.setLocation(loc.AppendLocation("examples")); err != nil {
		return err
	}
	if err := mt.Encoding.setLocation(loc.AppendLocation("encoding")); err != nil {
		return err
	}

	return nil
}
func (*MediaType) Kind() Kind      { return KindMediaType }
func (*MediaType) mapKind() Kind   { return KindMediaTypeMap }
func (*MediaType) sliceKind() Kind { return KindUndefined }

func (mt *MediaType) isNil() bool { return mt == nil }

var _ node = (*MediaType)(nil)
