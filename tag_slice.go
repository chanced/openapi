package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonpointer"
)

type TagSlice struct {
	Location `json:"-"`

	Items []Tag
}

// Kind implements node
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

// ResolveNodeByPointer implements node
func (ts TagSlice) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return ts.resolveNodeByPointer(ptr)
}

func (ts *TagSlice) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return ts, nil
	}
	nxt, tok, _ := ptr.Next()
	idx, err := tok.Int()
	if err != nil {
		return nil, newErrNotResolvable(ts.Location.AbsoluteLocation(), tok)
	}
	if idx < 0 {
		return nil, newErrNotResolvable(ts.Location.AbsoluteLocation(), tok)
	}
	if idx >= len(ts.Items) {
		return nil, newErrNotFound(ts.Location.AbsoluteLocation(), tok)
	}
	return ts.Items[idx].resolveNodeByPointer(nxt)
}

// MarshalJSON implements node
func (ts *TagSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.Items)
}

// UnmarshalJSON implements node
func (ts *TagSlice) UnmarshalJSON(data []byte) error {
	var items []Tag
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*ts = TagSlice{
		Items: items,
	}
	return nil
}

// isNil implements node
func (ts *TagSlice) isNil() bool {
	return ts == nil
}

// location implements node
func (ts TagSlice) location() Location {
	return ts.Location
}

func (*TagSlice) mapKind() Kind   { return KindUndefined }
func (*TagSlice) sliceKind() Kind { return KindUndefined }

func (ts *TagSlice) setLocation(loc Location) error {
	ts.Location = loc
	return nil
}

func (TagSlice) Anchors() (*Anchors, error) { return nil, nil }

var (
	_ node   = (*TagSlice)(nil)
	_ Walker = (*TagSlice)(nil)
)
