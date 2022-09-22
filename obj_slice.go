package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

type ObjSlice[T node] struct {
	Location `json:"-"`
	Items    []T `json:"-"`
}

// Anchors implements node
func (os *ObjSlice[T]) Anchors() (*Anchors, error) {
	if os == nil {
		return nil, nil
	}
	var a *Anchors
	var err error
	for _, x := range os.Items {
		if a, err = a.merge(x.Anchors()); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// Kind implements node
func (os *ObjSlice[T]) Kind() Kind {
	var t T
	return objSliceKind(t)
}

func (os *ObjSlice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.Items)
}

// Refs implements node
func (os *ObjSlice[T]) Refs() []Ref {
	if os == nil {
		return nil
	}
	var refs []Ref
	for _, x := range os.Items {
		refs = append(refs, x.Refs()...)
	}
	return refs
}

// // ResolveNodeByPointer implements node
// func (os *ObjSlice[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return os.resolveNodeByPointer(ptr)
// }

// func (os *ObjSlice[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return os, nil
// 	}
// 	nxt, tok, _ := ptr.Next()

// 	idx, err := tok.Int()
// 	if err != nil || idx < 0 {
// 		return nil, newErrNotResolvable(os.absolute, tok)
// 	}
// 	if idx >= len(os.Items) {
// 		return nil, newErrNotFound(os.AbsoluteLocation(), tok)
// 	}
// 	return os.Items[idx].resolveNodeByPointer(nxt)
// }

// UnmarshalJSON implements node
func (os *ObjSlice[T]) UnmarshalJSON(data []byte) error {
	*os = ObjSlice[T]{}
	items := []T{}
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	os.Items = items
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (os ObjSlice[T]) MarshalYAML() (interface{}, error) {
	j, err := os.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (os *ObjSlice[T]) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, os)
}

func (os *ObjSlice[T]) nodes() []node {
	var edges []node
	for _, x := range os.Items {
		edges = appendEdges(edges, x)
	}
	return edges
}

func (os *ObjSlice[T]) setLocation(loc Location) error {
	if os == nil {
		return nil
	}
	os.Location = loc
	return nil
}

func (os *ObjSlice[T]) isNil() bool { return os == nil }

func (*ObjSlice[T]) sliceKind() Kind { return KindUndefined }
func (*ObjSlice[T]) mapKind() Kind   { return KindUndefined }

var _ node = (*ObjSlice[*SecurityRequirement])(nil)
