package openapi

import (
	"encoding/json"
	"strconv"

	"github.com/chanced/jsonpointer"
)

type SchemaSlice struct {
	Location
	Items []*Schema
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

func (*SchemaSlice) Kind() Kind { return KindSchemaSlice }

func (ss *SchemaSlice) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}

	return ss.resolveNodeByPointer(ptr)
}

func (ss *SchemaSlice) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return ss, nil
	}
	nxt, tok, _ := ptr.Next()
	idx, err := tok.Int()
	if err != nil {
		return nil, newErrNotResolvable(ss.Location.AbsoluteLocation(), tok)
	}
	if idx < 0 {
		return nil, newErrNotFound(ss.Location.AbsoluteLocation(), tok)
	}
	if idx >= len(ss.Items) {
		return nil, newErrNotFound(ss.Location.AbsoluteLocation(), tok)
	}
	return ss.Items[idx].resolveNodeByPointer(nxt)
}

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

func (*SchemaSlice) mapKind() Kind   { return KindUndefined }
func (*SchemaSlice) sliceKind() Kind { return KindUndefined }

func (ss *SchemaSlice) setLocation(loc Location) error {
	if ss == nil {
		return nil
	}
	ss.Location = loc
	for i, s := range ss.Items {
		err := s.setLocation(loc.Append(strconv.Itoa(i)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *SchemaSlice) isNil() bool { return ss == nil }

var _ node = (*SchemaSlice)(nil)
