package openapi

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

// Referencable is any object type which could also be a Reference
type Referencable interface {
	IsRef() bool
}

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
	Ref string `yaml:"$ref" json:"$ref"`
	// A short summary which by default SHOULD override that of the referenced
	// component. If the referenced object-type does not allow a summary field,
	// then this field has no effect.
	Summary string `yaml:"summary,omitempty" json:"summary,omitempty"`
	// A description which by default SHOULD override that of the referenced
	// component. CommonMark syntax MAY be used for rich text representation. If
	// the referenced object-type does not allow a description field, then this
	// field has no effect.
	Description string `yaml:"description" json:"description,omitempty"`
}

// ParameterType returns ParameterTypeReference
func (r *Reference) ParameterType() ParameterType {
	return ParameterTypeReference
}

// ResponseKind distinguishes Reference by returning HeaderKindRef
func (r *Reference) ResponseKind() ResponseKind {
	return ResponseKindRef
}

// ExampleKind distinguishes Reference by returning HeaderKindRef
func (r *Reference) ExampleKind() ExampleKind {
	return ExampleKindRef
}

// HeaderKind distinguishes Reference by returning HeaderKindRef
func (r *Reference) HeaderKind() HeaderKind {
	return HeaderKindRef
}

// RequestBodyKind returns RequestBodyKindRef
func (r *Reference) RequestBodyKind() RequestBodyKind {
	return RequestBodyKindRef
}

// CallbackKind returns CallbackKindRef
func (r *Reference) CallbackKind() CallbackKind {
	return CallbackKindRef
}

// PathKind returns PathKindRef
func (r *Reference) PathKind() PathKind {
	return PathKindRef
}

// SecuritySchemeKind returns SecuritySchemeKindRef
func (r *Reference) SecuritySchemeKind() SecuritySchemeKind {
	return SecuritySchemeKindRef
}

// LinkKind returns LinkKindRef
func (r *Reference) LinkKind() LinkKind {
	return LinkKindRef
}

func isRefJSON(data []byte) bool {
	r := gjson.GetBytes(data, "$ref")
	return r.Str != ""
}

func unmarshalReferenceJSON(data []byte) (*Reference, error) {
	if !isRefJSON(data) {
		return nil, ErrNotReference
	}
	var r Reference
	return &r, json.Unmarshal(data, &r)
}

var _ SecurityScheme = (*Reference)(nil)
var _ Path = (*Reference)(nil)
var _ Response = (*Reference)(nil)
var _ Example = (*Reference)(nil)
var _ Parameter = (*Reference)(nil)
var _ Header = (*Reference)(nil)
