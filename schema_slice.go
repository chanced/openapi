package openapi

import (
	"encoding/json"
	"strconv"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

type SchemaSlice struct {
	Location
	Items []*Schema
}

func (ss *SchemaSlice) Nodes() []Node {
	if ss == nil {
		return nil
	}
	return downcastNodes(ss.nodes())
}

func (ss *SchemaSlice) nodes() []node {
	edges := make([]node, len(ss.Items))
	for i, s := range ss.Items {
		edges[i] = s
	}
	return edges
}

func (ss *SchemaSlice) Anchors() (*Anchors, error) {
	if ss == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	for _, s := range ss.Items {
		if anchors, err = anchors.merge(s.Anchors()); err != nil {
			return nil, err
		}
	}
	return anchors, nil
}

func (ss *SchemaSlice) Refs() []Ref {
	if ss == nil {
		return nil
	}
	var refs []Ref
	for _, s := range ss.Items {
		refs = append(refs, s.Refs()...)
	}
	return refs
}

func (*SchemaSlice) Kind() Kind { return KindSchemaSlice }

// func (ss *SchemaSlice) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}

// 	return ss.resolveNodeByPointer(ptr)
// }

// func (ss *SchemaSlice) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return ss, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	idx, err := tok.Int()
// 	if err != nil {
// 		return nil, newErrNotResolvable(ss.Location.AbsoluteLocation(), tok)
// 	}
// 	if idx < 0 {
// 		return nil, newErrNotFound(ss.Location.AbsoluteLocation(), tok)
// 	}
// 	if idx >= len(ss.Items) {
// 		return nil, newErrNotFound(ss.Location.AbsoluteLocation(), tok)
// 	}
// 	return ss.Items[idx].resolveNodeByPointer(nxt)
// }

func (ss SchemaSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(ss.Items)
}

func (ss *SchemaSlice) UnmarshalJSON(data []byte) error {
	var v []*Schema
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*ss = SchemaSlice{Items: v}
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (ss SchemaSlice) MarshalYAML() (interface{}, error) {
	j, err := ss.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(j, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (ss *SchemaSlice) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, ss)
}

func (*SchemaSlice) mapKind() Kind   { return KindUndefined }
func (*SchemaSlice) sliceKind() Kind { return KindUndefined }

func (ss *SchemaSlice) setLocation(loc Location) error {
	if ss == nil {
		return nil
	}
	ss.Location = loc
	for i, s := range ss.Items {
		err := s.setLocation(loc.AppendLocation(strconv.Itoa(i)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *SchemaSlice) Clone() *SchemaSlice {
	if ss == nil {
		return nil
	}
	v := make([]*Schema, len(ss.Items))
	for i, s := range ss.Items {
		v[i] = s.Clone()
	}
	return &SchemaSlice{Items: v}
}
func (ss *SchemaSlice) isNil() bool { return ss == nil }

var _ node = (*SchemaSlice)(nil)
