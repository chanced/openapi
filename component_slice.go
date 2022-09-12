package openapi

import (
	"encoding/json"
	"strconv"

	"github.com/chanced/jsonpointer"
)

// ComponentSlice is a slice of Components of type T
type ComponentSlice[T node] struct {
	Location `json:"-"`
	Items    []Component[T] `json:"-"`
}

func (ComponentSlice[T]) Kind() Kind {
	var t T
	return t.Kind()
}

func (cs ComponentSlice[T]) Resolve(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}

	return cs.resolve(ptr)
}

func (cs *ComponentSlice[T]) resolve(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return cs, nil
	}
	nxt, tok, _ := ptr.Next()
	idx, err := tok.Int()
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(cs.Items) {
		return nil, newErrNotFound(cs.Location.AbsoluteLocation(), tok)
	}
	return cs.Items[idx].resolve(nxt)
}

func (cs ComponentSlice[T]) MarshalJSON() ([]byte, error) {
	type componentslice[T node] ComponentSlice[T]
	return json.Marshal(componentslice[T](cs))
}

func (cs *ComponentSlice[T]) UnmarshalJSON(data []byte) error {
	type componentslice[T node] ComponentSlice[T]
	var v componentslice[T]
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*cs = ComponentSlice[T](v)
	return nil
}

func (*ComponentSlice[T]) mapKind() Kind { return KindUndefined }

func (ComponentSlice[T]) sliceKind() Kind {
	var t T
	return t.sliceKind()
}

func (cs *ComponentSlice[T]) setLocation(loc Location) error {
	cs.Location = loc
	for i, c := range cs.Items {
		if err := c.setLocation(loc.Append(strconv.Itoa(i))); err != nil {
			return err
		}
	}
	return nil
}

var _ node = (*ComponentSlice[*Server])(nil)
