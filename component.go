package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

type Component[T refable] struct {
	Location
	Reference *Reference[T]
	Object    T
}

func (c *Component[T]) nodes() []node {
	if c == nil {
		return nil
	}
	if c.IsReference() {
		return appendEdges(nil, c.Reference)
	}
	return appendEdges(nil, c.Object)
}

// Edges returns the immediate edges of the Node. This is used to build a
// graph of the OpenAPI document.
//

// IsResolved implements Ref
func (c *Component[T]) IsResolved() bool {
	if c == nil || c.Reference == nil {
		return true
	}
	return !c.Object.isNil()
}

// URI implements Ref
func (c *Component[T]) URI() *uri.URI {
	if !c.IsReference() {
		return nil
	}
	return c.Reference.Ref
}

// MakeReference converts the Component into a reference, altering the path of
// all nested nodes.
func (c *Component[T]) MakeReference(ref uri.URI) error {
	if c.Object.isNil() {
		return fmt.Errorf("cannot make reference to nil object")
	}
	c.Reference.dst = &c.Object
	c.Reference.Ref = &ref
	loc, err := NewLocation(ref)
	if err != nil {
		return err
	}

	return c.Object.setLocation(loc)
}

func (c *Component[T]) Kind() Kind {
	switch c.ObjectKind() {
	case KindExample:
		return KindExampleComponent
	case KindHeader:
		return KindHeaderComponent
	case KindServer:
		return KindServerComponent
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

// IsRef returns false
//

// IsReference returns true if this Component contains a Reference
func (c *Component[T]) IsReference() bool { return !c.Reference.isNil() }

func (c *Component[T]) Refs() []Ref {
	if c == nil {
		return nil
	}
	if c.IsReference() {
		return []Ref{c.Reference}
	}
	return c.Object.Refs()
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

// ComponentKind returns the Kind of the containing Object, regardless of if it
// is referenced or not.
func (c *Component[T]) ObjectKind() Kind {
	return c.Object.Kind()
}

func (c *Component[T]) UnmarshalJSON(data []byte) error {
	if isRefJSON(data) {
		var ref Reference[T]
		if err := json.Unmarshal(data, &ref); err != nil {
			return err
		}

		ref.ReferencedKind = c.ObjectKind()
		ref.dst = &c.Object
		c.Reference = &ref

		*c = Component[T]{
			Reference: &ref,
		}
		return nil
	}
	var obj T

	k := obj.Kind()
	_ = k

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*c = Component[T]{
		Object: obj,
	}
	return nil
}

func (c *Component[T]) MarshalYAML() (interface{}, error) {
	j, err := c.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Component[T]) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, c)
}

func (c *Component[T]) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = loc
	if c.Reference != nil {
		return c.Reference.setLocation(loc)
	} else if !c.Object.isNil() {
		return c.Object.setLocation(loc)
	}
	return nil
}

func (c *Component[T]) Anchors() (*Anchors, error) {
	if c == nil {
		return nil, nil
	}
	if c.Reference != nil {
		return nil, nil
	}
	return c.Object.Anchors()
}

func (c *Component[T]) isNil() bool { return c == nil }

var _ node = (*Component[*Response])(nil)

// func (c *Component[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return c.resolveNodeByPointer(ptr)
// }

// func (c *Component[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return c, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "$ref":
// 		if nxt.IsRoot() {
// 			return c.Reference, nil
// 		}
// 		return nil, newErrNotResolvable(c.Location.AbsoluteLocation(), tok)
// 	default:
// 		// TODO: this may need to change. Not sure when I need to perform these
// 		// resolutions just yet. If before population, Object may be nil at this call.
// 		return c.Object.resolveNodeByPointer(nxt)
// 	}
// }
