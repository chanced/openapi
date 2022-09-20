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
type Reference struct {
	// The reference identifier. This MUST be in the form of a URI.
	//
	// 	*required*
	Ref *uri.URI `yaml:"$ref" json:"$ref"`

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

	dst interface{}

	resolved bool
}

func (r *Reference) Nodes() []Node {
	if r == nil {
		return nil
	}
	return downcastNodes(r.nodes())
}

func (r *Reference) nodes() []node {
	if r == nil {
		return nil
	}

	return appendEdges(nil, r.Resolved().(node))
}
func (r *Reference) RefKind() Kind { return r.ReferencedKind }

func (r *Reference) URI() *uri.URI {
	if r == nil {
		return nil
	}
	return r.Ref
}

func (r *Reference) IsResolved() bool { return r.resolved }

// resolve resolves the reference
//
// TODO: make this a bit less panicky
func (r *Reference) resolve(v Node) error {
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
		rv.Elem().Set(reflect.ValueOf(v).Elem())
	}
	return nil
}

// Referenced returns the resolved referenced Node
func (r *Reference) Resolved() Node {
	n, ok := (r.dst).(*Node)
	if !ok {
		panic("openapi: Reference dst is not a Node. This is a bug. Please report it: https://github.com/chanced/openapi/issues/new")
	}
	return *n
}

// Refs returns nil as instances of Reference do not contain the referenced
// object
func (*Reference) Refs() []Ref { return nil }

func (r *Reference) Anchors() (*Anchors, error) { return nil, nil }

func (*Reference) RefType() RefType { return RefTypeComponent }

// func (r *Reference) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return r.resolveNodeByPointer(ptr)
// }

// func (r *Reference) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return r, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(r.Location.AbsoluteLocation(), tok)
// }

func (r Reference) MarshalJSON() ([]byte, error) {
	type reference Reference
	return json.Marshal(reference(r))
}

func (r *Reference) UnmarshalJSON(data []byte) error {
	type reference Reference
	var v reference
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Reference(v)
	return nil
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (r Reference) MarshalYAML() (interface{}, error) {
	j, err := r.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (r *Reference) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, r)
}

func (r *Reference) String() string {
	return r.Ref.String()
}

func (r *Reference) setLocation(loc Location) error {
	if r == nil {
		return nil
	}
	r.Location = loc
	return nil
}

func (r *Reference) Kind() Kind    { return KindReference }
func (*Reference) mapKind() Kind   { return KindUndefined }
func (*Reference) sliceKind() Kind { return KindUndefined }

func (r *Reference) isNil() bool { return r == nil }

func isRefJSON(data []byte) bool {
	r := gjson.GetBytes(data, "$ref")
	return r.Str != ""
}

var (
	_ node = (*Reference)(nil)

	_ Ref = (*Reference)(nil)
	_ ref = (*Reference)(nil)
)
