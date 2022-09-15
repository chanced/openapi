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

// Refs implements node
func (cs *ComponentSlice[T]) Refs() []Ref {
	var refs []Ref
	for _, item := range cs.Items {
		refs = append(refs, item.Refs()...)
	}
	return refs
}

func (ComponentSlice[T]) Kind() Kind {
	var t T
	return t.Kind()
}

func (cs ComponentSlice[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}

	return cs.resolveNodeByPointer(ptr)
}

func (cs *ComponentSlice[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
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
	return cs.Items[idx].resolveNodeByPointer(nxt)
}

func (cs ComponentSlice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs.Items)
}

func (cs *ComponentSlice[T]) UnmarshalJSON(data []byte) error {
	var items []Component[T]
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*cs = ComponentSlice[T]{
		Items: items,
	}
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

func (cs *ComponentSlice[T]) Anchors() (*Anchors, error) {
	if cs == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	for _, item := range cs.Items {
		if anchors, err = item.Anchors(); err != nil {
			return nil, err
		}
	}
	return anchors, nil
}

func (cs *ComponentSlice[T]) isNil() bool { return cs == nil }

var (
	_ node   = (*ComponentSlice[*Server])(nil)
	_ Walker = (*ComponentSlice[*Server])(nil)
)

// func (cs *ComponentSlice[T]) Walk(v Visitor) error {
// 	var t T
// 	var err error
// 	v, err = v.Visit(cs)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}
// 	switch t.Kind() {
// 	case KindParameter:
// 		return cs.walkParameters(v)
// 	case KindServer:
// 		return cs.walkServers(v)
// 	default:

// 	}
// }

// func (cs *ComponentSlice[T]) walkParameters(v Visitor) error {
// 	var err error
// 	ps, ok := (any)(cs).(*ComponentSlice[*Parameter])
// 	if !ok {
// 		// shouldn't happen
// 		panic(fmt.Sprintf("%T is not a *ComponentSlice[*Parameter]", cs))
// 	}
// 	v, err = v.VisitParameterSlice(ps)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}
// 	for _, p := range ps.Items {
// 		if err = p.Walk(v); err != nil {
// 			return err
// 		}
// 	}
// }
