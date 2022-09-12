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

func (*SchemaSlice) Kind() Kind { return KindSchemaSlice }

func (ss *SchemaSlice) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}

	return ss.resolve(ptr)
}

func (ss *SchemaSlice) resolve(ptr jsonpointer.Pointer) (Node, error) {
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
	return ss.Items[idx].resolve(nxt)
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

var _ node = (*SchemaSlice)(nil)
