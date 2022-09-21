package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

type TagSlice struct {
	Location `json:"-"`

	Items []*Tag
}

func (tl *TagSlice) Nodes() []Node {
	if tl == nil {
		return nil
	}
	return downcastNodes(tl.nodes())
}

func (tl *TagSlice) nodes() []node {
	edges := make([]node, len(tl.Items))
	for i, item := range tl.Items {
		edges[i] = item
	}
	return edges
}

func (*TagSlice) Kind() Kind { return KindTagSlice }

func (ts *TagSlice) Refs() []Ref {
	if ts == nil {
		return nil
	}
	var refs []Ref
	for _, item := range ts.Items {
		refs = append(refs, item.Refs()...)
	}
	return refs
}

// func (ts TagSlice) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return ts.resolveNodeByPointer(ptr)
// }

// func (ts *TagSlice) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return ts, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	idx, err := tok.Int()
// 	if err != nil {
// 		return nil, newErrNotResolvable(ts.Location.AbsoluteLocation(), tok)
// 	}
// 	if idx < 0 {
// 		return nil, newErrNotResolvable(ts.Location.AbsoluteLocation(), tok)
// 	}
// 	if idx >= len(ts.Items) {
// 		return nil, newErrNotFound(ts.Location.AbsoluteLocation(), tok)
// 	}
// 	return ts.Items[idx].resolveNodeByPointer(nxt)
// }

func (ts *TagSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.Items)
}

func (ts *TagSlice) UnmarshalJSON(data []byte) error {
	var items []*Tag
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*ts = TagSlice{
		Items: items,
	}
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (t TagSlice) MarshalYAML() (interface{}, error) {
	j, err := t.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (t *TagSlice) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, t)
}

func (ts *TagSlice) isNil() bool {
	return ts == nil
}

func (ts TagSlice) location() Location {
	return ts.Location
}

func (*TagSlice) mapKind() Kind   { return KindUndefined }
func (*TagSlice) sliceKind() Kind { return KindUndefined }

func (ts *TagSlice) setLocation(loc Location) error {
	if ts == nil {
		return nil
	}
	ts.Location = loc
	return nil
}

func (*TagSlice) Anchors() (*Anchors, error) { return nil, nil }

var _ node = (*TagSlice)(nil)
