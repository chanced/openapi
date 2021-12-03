package openapi

import (
	"encoding/json"
	"errors"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
)

type ParameterResolver func(ref string) (*ParameterObj, error)

type ResponseResolver func(ref string) (*ResponseObj, error)

type ExampleResolver func(ref string) (*ExampleObj, error)

type HeaderResolver func(ref string) (*HeaderObj, error)

type RequestBodyResolver func(ref string) (*RequestBodyObj, error)

type CallbackResolver func(ref string) (*CallbackObj, error)

type PathResolver func(ref string) (*PathObj, error)

type SecuritySchemeResolver func(ref string) (*SecuritySchemeObj, error)

type LinkResolver func(ref string) (*LinkObj, error)

type SchemaResolver func(ref string) (*SchemaObj, error)

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

// MarshalYAML marshals YAML
func (r Reference) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(r)
}

// UnmarshalYAML unmarshals YAML
func (r *Reference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, r)
}

// ParameterKind returns ParameterKindReference
func (r *Reference) ParameterKind() ParameterKind {
	return ParameterKindReference
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

// ResolveParameter resolves r by invoking resolve
func (r *Reference) ResolveParameter(resolve ParameterResolver) (*ParameterObj, error) {
	return resolve(r.Ref)
}

// ResolveResponse resolves r by invoking resolve
func (r *Reference) ResolveResponse(resolve ResponseResolver) (*ResponseObj, error) {
	return resolve(r.Ref)
}

// ResolveExample resolves r by invoking resolve
func (r *Reference) ResolveExample(resolve ExampleResolver) (*ExampleObj, error) {
	return resolve(r.Ref)
}

// ResolveHeader resolves r by invoking resolve
func (r *Reference) ResolveHeader(resolve HeaderResolver) (*HeaderObj, error) {
	return resolve(r.Ref)
}

// ResolveRequestBody resolves r by invoking resolve
func (r *Reference) ResolveRequestBody(resolve RequestBodyResolver) (*RequestBodyObj, error) {
	return resolve(r.Ref)
}

// ResolveCallback resolves r by invoking resolve
func (r *Reference) ResolveCallback(resolve CallbackResolver) (*CallbackObj, error) {
	return resolve(r.Ref)
}

// ResolvePath resolves r by invoking resolve
func (r *Reference) ResolvePath(resolve PathResolver) (*PathObj, error) {
	return resolve(r.Ref)
}

// ResolveSecurityScheme resolves r by invoking resolve
func (r *Reference) ResolveSecurityScheme(resolve SecuritySchemeResolver) (*SecuritySchemeObj, error) {
	return resolve(r.Ref)
}

// ResolveLink resolves r by invoking resolve
func (r *Reference) ResolveLink(resolve LinkResolver) (*LinkObj, error) {
	return resolve(r.Ref)
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
