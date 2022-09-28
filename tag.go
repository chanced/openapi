package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

type TagMap = ObjMap[*Tag]

// Tag adds metadata that is used by the Operation Object.
//
// It is not mandatory to have a Tag Object per tag defined in the Operation
// Object instances.
type Tag struct {
	Location   `json:"-"`
	Extensions `json:"-"`
	// The name of the tag.
	//
	// 	*required*
	Name Text `json:"name"`
	//  A description for the tag.
	//
	// CommonMark syntax MAY be used for rich text representation.
	//
	// https://spec.commonmark.org/
	Description Text `json:"description,omitempty"`
	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" bson:"externalDocs,omitempty"`
}

func (t *Tag) Nodes() []Node {
	if t == nil {
		return nil
	}
	return downcastNodes(t.nodes())
}

func (t *Tag) nodes() []node {
	if t == nil {
		return nil
	}
	return appendEdges(nil, t.ExternalDocs)
}

func (*Tag) Refs() []Ref                { return nil }
func (*Tag) Anchors() (*Anchors, error) { return nil, nil }
func (t *Tag) isNil() bool              { return t == nil }

// func (t *Tag) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return t.resolveNodeByPointer(ptr)
// }

// func (t *Tag) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return t, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "externalDocs":
// 		if nxt.IsRoot() {
// 			return t.ExternalDocs, nil
// 		}
// 		if t.ExternalDocs == nil {
// 			return nil, newErrNotFound(t.Location.AbsoluteLocation(), tok)
// 		}
// 		return t.ExternalDocs.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(t.Location.AbsoluteLocation(), tok)
// 	}
// }

func (*Tag) Kind() Kind      { return KindTag }
func (*Tag) mapKind() Kind   { return KindUndefined }
func (*Tag) sliceKind() Kind { return KindTagSlice }

func (t *Tag) setLocation(loc Location) error {
	if t == nil {
		return nil
	}
	t.Location = loc
	return t.ExternalDocs.setLocation(loc.AppendLocation("externalDocs"))
}

// MarshalJSON marshals t into JSON
func (t Tag) MarshalJSON() ([]byte, error) {
	type tag Tag

	return marshalExtendedJSON(tag(t))
}

// UnmarshalJSON unmarshals json into t
func (t *Tag) UnmarshalJSON(data []byte) error {
	type tag Tag

	v := tag{}
	err := unmarshalExtendedJSON(data, &v)
	*t = Tag(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (t Tag) MarshalYAML() (interface{}, error) {
	j, err := t.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (t *Tag) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, t)
}

var _ node = (*Tag)(nil)
