package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// ComponentMap is a pseudo map consisting of Components with type T.

func newComponent[T refable](ref *Reference[T], obj T) Component[T] {
	return Component[T]{
		Reference: ref,
		Object:    obj,
	}
}

// Components holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the
// API unless they are explicitly referenced from properties outside the
// components object.
type Components struct {
	// OpenAPI extensions
	Extensions `json:"-"`
	Location   `json:"-"`

	Schemas         *SchemaMap         `json:"schemas,omitempty"`
	Responses       *ResponseMap       `json:"responses,omitempty"`
	Parameters      *ParameterMap      `json:"parameters,omitempty"`
	RequestBodies   *RequestBodyMap    `json:"requestBodies,omitempty"`
	Headers         *HeaderMap         `json:"headers,omitempty"`
	SecuritySchemes *SecuritySchemeMap `json:"securitySchemes,omitempty"`
	Links           *LinkMap           `json:"links,omitempty"`
	Callbacks       *CallbacksMap      `json:"callbacks,omitempty"`
	PathItems       *PathItemMap       `json:"pathItems,omitempty"`
	Examples        *ExampleMap        `json:"examples,omitempty"` //
}

func (*Components) Kind() Kind { return KindComponents }

func (c *Components) Refs() []Ref {
	if c == nil {
		return nil
	}
	var refs []Ref
	if c.Schemas != nil {
		refs = append(refs, c.Schemas.Refs()...)
	}
	if c.Responses != nil {
		refs = append(refs, c.Responses.Refs()...)
	}
	if c.Parameters != nil {
		refs = append(refs, c.Parameters.Refs()...)
	}
	if c.Examples != nil {
		refs = append(refs, c.Examples.Refs()...)
	}
	if c.RequestBodies != nil {
		refs = append(refs, c.RequestBodies.Refs()...)
	}
	if c.Headers != nil {
		refs = append(refs, c.Headers.Refs()...)
	}
	if c.SecuritySchemes != nil {
		refs = append(refs, c.SecuritySchemes.Refs()...)
	}
	if c.Links != nil {
		refs = append(refs, c.Links.Refs()...)
	}
	if c.Callbacks != nil {
		refs = append(refs, c.Callbacks.Refs()...)
	}
	if c.PathItems != nil {
		refs = append(refs, c.PathItems.Refs()...)
	}
	return refs
}

// func (c *Components) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	if c == nil {
// 		return nil, nil
// 	}
// 	return c.resolveNodeByPointer(ptr)
// }

//	func (c *Components) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
//		if ptr.IsRoot() {
//			return c, nil
//		}
//		nxt, tok, _ := ptr.Next()
//		switch tok {
//		case "schemas":
//			if c.Schemas == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Schemas.resolveNodeByPointer(nxt)
//		case "responses":
//			if c.Responses == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Responses.resolveNodeByPointer(nxt)
//		case "parameters":
//			if c.Parameters == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Parameters.resolveNodeByPointer(nxt)
//		case "examples":
//			if c.Examples == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Examples.resolveNodeByPointer(nxt)
//		case "requestBodies":
//			if c.RequestBodies == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.RequestBodies.resolveNodeByPointer(nxt)
//		case "headers":
//			if c.Headers == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Headers.resolveNodeByPointer(nxt)
//		case "securitySchemes":
//			if c.SecuritySchemes == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.SecuritySchemes.resolveNodeByPointer(nxt)
//		case "links":
//			if c.Links == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Links.resolveNodeByPointer(nxt)
//		case "callbacks":
//			if c.Callbacks == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.Callbacks.resolveNodeByPointer(nxt)
//		case "pathItems":
//			if c.PathItems == nil {
//				return nil, newErrNotFound(c.AbsoluteLocation(), tok)
//			}
//			return c.PathItems.resolveNodeByPointer(nxt)
//		default:
//			return nil, newErrNotResolvable(c.AbsoluteLocation(), tok)
//		}
//	}
func (c *Components) nodes() []node {
	if c == nil {
		return nil
	}
	edges := appendEdges(nil, c.Schemas)
	edges = appendEdges(edges, c.Responses)
	edges = appendEdges(edges, c.Parameters)
	edges = appendEdges(edges, c.Examples)
	edges = appendEdges(edges, c.RequestBodies)
	edges = appendEdges(edges, c.Headers)
	edges = appendEdges(edges, c.SecuritySchemes)
	edges = appendEdges(edges, c.Links)
	edges = appendEdges(edges, c.Callbacks)
	edges = appendEdges(edges, c.PathItems)
	return edges
}

func (c *Components) isNil() bool {
	return c == nil
}

func (*Components) mapKind() Kind   { return KindUndefined }
func (*Components) sliceKind() Kind { return KindUndefined }

func (c *Components) Anchors() (*Anchors, error) {
	if c == nil {
		return nil, nil
	}
	var err error
	var anchors *Anchors
	if anchors, err = anchors.merge(c.Schemas.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Responses.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Parameters.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Examples.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.RequestBodies.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Headers.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.SecuritySchemes.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Links.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.Callbacks.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(c.PathItems.Anchors()); err != nil {
		return nil, err
	}

	return anchors, nil
}

func (c *Components) setLocation(loc Location) error {
	if c == nil {
		return nil
	}
	c.Location = loc
	var err error
	if err = c.Schemas.setLocation(loc.AppendLocation("schemas")); err != nil {
		return err
	}
	if err = c.Responses.setLocation(loc.AppendLocation("responses")); err != nil {
		return err
	}
	if err = c.Parameters.setLocation(loc.AppendLocation("parameters")); err != nil {
		return err
	}
	if err = c.Examples.setLocation(loc.AppendLocation("examples")); err != nil {
		return err
	}
	if err = c.RequestBodies.setLocation(loc.AppendLocation("requestBodies")); err != nil {
		return err
	}
	if err = c.Headers.setLocation(loc.AppendLocation("headers")); err != nil {
		return err
	}
	if err = c.SecuritySchemes.setLocation(loc.AppendLocation("securitySchemes")); err != nil {
		return err
	}
	if err = c.Links.setLocation(loc.AppendLocation("links")); err != nil {
		return err
	}
	if err = c.Callbacks.setLocation(loc.AppendLocation("callbacks")); err != nil {
		return err
	}
	if err = c.PathItems.setLocation(loc.AppendLocation("pathItems")); err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals JSON
func (c Components) MarshalJSON() ([]byte, error) {
	type components Components
	return marshalExtendedJSON(components(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Components) UnmarshalJSON(data []byte) error {
	type components Components
	var v components
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*c = Components(v)
	return nil
}

func (c Components) MarshalYAML() (interface{}, error) {
	j, err := c.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Components) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, c)
}

var _ node = (*Components)(nil)
