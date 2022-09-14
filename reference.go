package openapi

import (
	"encoding/json"
	"errors"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
	"github.com/tidwall/gjson"
)

// ErrNotReference indicates not a reference
var ErrNotReference = errors.New("error: data is not a Reference")

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
}

func (r *Reference) URI() *uri.URI {
	return r.Ref
}

func (r *Reference) Anchors() (*Anchors, error) { return nil, nil }

func (r *Reference) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return r.resolveNodeByPointer(ptr)
}

func (r *Reference) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return r, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(r.Location.AbsoluteLocation(), tok)
}

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

var _ node = (*Reference)(nil)
