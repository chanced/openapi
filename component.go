package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type Component[T node] struct {
	Location
	Reference *Reference
	Object    T
}

// Dst implements Ref
func (c *Component[T]) RefDst() interface{} {
	if !c.IsRef() || c == nil {
		return nil
	}
	// this should always be a pointer
	return &c.Object
}

// IsResolved implements Ref
func (c *Component[T]) IsResolved() bool {
	if c == nil || c.Reference == nil {
		return true
	}
	return !c.Object.isNil()
}

// RefURI implements Ref
func (c *Component[T]) RefURI() *uri.URI {
	if !c.IsRef() {
		return nil
	}
	return c.Reference.Ref
}

func (c *Component[T]) Kind() Kind {
	switch c.Object.Kind() {
	case KindExample:
		return KindExampleComponent
	case KindHeader:
		return KindHeaderComponent
	case KindLink:
		return KindLinkComponent
	case KindResponse:
		return KindResponseComponent
	case KindParameter:
		return KindParameterComponent
	case KindPathItem:
		return KindPathItemComponent
	case KindRequestBody:
		return KindRequestBodyComponent
	case KindCallbacks:
		return KindCallbacksComponent
	case KindSecurityScheme:
		return KindSecuritySchemeComponent
	default:
		return KindUndefined
	}
}

func (c *Component[T]) location() Location {
	if c.Reference != nil {
		return c.Reference.Location
	}
	return c.Object.location()
}

func (c *Component[T]) IsRef() bool {
	return c.Reference != nil
}

func (c *Component[T]) Refs() []Ref {
	if c == nil {
		return nil
	}
	if c.IsRef() {
		return []Ref{c}
	}
	return c.Object.Refs()
}

// func (c *Component[T]) IsResolved() bool {

// }

func (c *Component[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return c.resolveNodeByPointer(ptr)
}

func (c *Component[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return c, nil
	}
	nxt, tok, _ := ptr.Next()
	switch tok {
	case "$ref":
		if nxt.IsRoot() {
			return c.Reference, nil
		}
		return nil, newErrNotResolvable(c.Location.AbsoluteLocation(), tok)
	default:
		// TODO: this may need to change. Not sure when I need to perform these
		// resolutions just yet. If before population, Object may be nil at this call.
		return c.Object.resolveNodeByPointer(nxt)
	}
}

func (*Component[T]) mapKind() Kind {
	var t T
	return t.mapKind()
}

func (*Component[T]) sliceKind() Kind {
	var t T
	return t.sliceKind()
}

func (c Component[T]) MarshalJSON() ([]byte, error) {
	if c.Reference != nil {
		return json.Marshal(c.Reference)
	}
	if any(c.Object) != nil {
		return c.Object.MarshalJSON()
	}
	return nil, nil
}

func (c *Component[T]) UnmarshalJSON(data []byte) error {
	if isRefJSON(data) {
		var ref Reference
		if err := json.Unmarshal(data, &ref); err != nil {
			return err
		}
		*c = Component[T]{
			Reference: &ref,
		}
		return nil
	}
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*c = Component[T]{
		Object: obj,
	}
	return nil
}

func (c *Component[T]) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	if c.Reference != nil {
		return c.Reference.setLocation(loc)
	} else if !c.Object.isNil() {
		return c.Object.setLocation(loc)
	}
	return nil
}

func (c *Component[T]) Anchors() (*Anchors, error) {
	if c.Reference != nil {
		return nil, nil
	}
	return c.Object.Anchors()
}

func (c *Component[T]) isNil() bool { return c == nil }

var (
	_ node   = (*Component[*Server])(nil)
	_ Walker = (*Component[*Server])(nil)
)
