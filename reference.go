package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

// ErrNotReference indicates not a reference
var ErrNotReference = errors.New("openapi: data is not a Reference")

// Reference is simple object to allow referencing other components in the
// OpenAPI document, internally and externally.
//
// The $ref string value contains a URI
// [RFC3986](https://datatracker.ietf.org/doc/html/rfc3986), which identifies
// the location of the value being referenced.
//
// See the [rules for resolving Relative
// References](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#relativeReferencesURI).
type Reference[T refable] struct {
	// The reference identifier. This MUST be in the form of a URI.
	//
	// 	*required*
	Ref *uri.URI `json:"$ref"`

	// A short summary which by default SHOULD override that of the referenced
	// component. If the referenced object-type does not allow a summary field,
	// then this field has no effect.
	Summary Text `json:"summary,omitempty"`

	// A description which by default SHOULD override that of the referenced
	// component. CommonMark syntax MAY be used for rich text representation. If
	// the referenced object-type does not allow a description field, then this
	// field has no effect.
	Description Text `json:"description,omitempty"`

	// Location of the Reference
	Location `json:"-"`

	ReferencedKind Kind `json:"-"`

	Resolved T `json:"-"`

	dst interface{}

	resolved bool
}

func (r *Reference[T]) Nodes() []Node {
	if r == nil {
		return nil
	}
	return downcastNodes(r.nodes())
}

func (r *Reference[T]) nodes() []node {
	if r == nil {
		return nil
	}

	return appendEdges(nil, r.ResolvedNode().(node))
}
func (r *Reference[T]) RefKind() Kind { return r.ReferencedKind }

func (r *Reference[T]) URI() *uri.URI {
	if r == nil {
		return nil
	}
	return r.Ref
}

func (r *Reference[T]) IsResolved() bool { return r.resolved }

// resolve resolves the reference
//
// TODO: make this a bit less panicky
func (r *Reference[T]) resolve(v Node) error {
	if r == nil {
		return fmt.Errorf("openapi: Reference is nil")
	}
	if r.dst == nil {
		return fmt.Errorf("openapi: Reference dst is nil")
	}
	if v.Kind() != r.ReferencedKind {
		return NewResolutionError(r, r.ReferencedKind, v.Kind())
	}

	rd := reflect.ValueOf(r.dst)
	rv := reflect.ValueOf(v)
	if rv.Type().AssignableTo(rd.Type().Elem()) {
		rd.Elem().Set(rv)
	} else {
		return fmt.Errorf("%s is not assignable to %s", rv.Type().String(), rd.Type().String())
	}

	r.Resolved = v.(T)

	return nil
}

// Referenced returns the resolved referenced Node
func (r *Reference[T]) ResolvedNode() Node {
	return r.Resolved
}

// Refs returns nil as instances of Reference do not contain the referenced
// object
func (*Reference[T]) Refs() []Ref { return nil }

func (r *Reference[T]) Anchors() (*Anchors, error) { return nil, nil }

func (*Reference[T]) RefType() RefType { return RefTypeComponent }

// func (r *Reference[T]) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return r.resolveNodeByPointer(ptr)
// }

// func (r *Reference[T]) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return r, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(r.Location.AbsoluteLocation(), tok)
// }

func (r Reference[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(reference[T](r))
}

type reference[T refable] Reference[T]

func (r *Reference[T]) UnmarshalJSON(data []byte) error {
	var v reference[T]
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Reference[T](v)
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (r Reference[T]) MarshalYAML() (interface{}, error) {
	j, err := r.MarshalJSON()
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
func (r *Reference[T]) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, r)
}

func (r *Reference[T]) String() string {
	return r.Ref.String()
}

func (r *Reference[T]) setLocation(loc Location) error {
	if r == nil {
		return nil
	}
	r.Location = loc
	return nil
}

func (r *Reference[T]) Kind() Kind    { return KindReference }
func (*Reference[T]) mapKind() Kind   { return KindUndefined }
func (*Reference[T]) sliceKind() Kind { return KindUndefined }

func (r *Reference[T]) isNil() bool { return r == nil }

func isRefJSON(data []byte) bool {
	r := gjson.GetBytes(data, "$ref")
	return r.Str != ""
}

var (
	_ node = (*Reference[*Response])(nil)
	_ Ref  = (*Reference[*Response])(nil)
	_ ref  = (*Reference[*Response])(nil)
)
