package openapi

import (
	"encoding/json"
	"errors"

	"github.com/chanced/openapi/yamlutil"
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

// Kind returns KindReference
func (r *Reference) ParameterKind() Kind {
	return KindReference
}

// ResolveParameter resolves r by invoking resolve
func (r *Reference) ResolveParameter(resolve func(ref string) (*ParameterObj, error)) (*ParameterObj, error) {
	return resolve(r.Ref)
}

// ResolveResponse resolves r by invoking resolve
func (r *Reference) ResolveResponse(resolve func(ref string) (*ResponseObj, error)) (*ResponseObj, error) {
	return resolve(r.Ref)
}

// ResolveExample resolves r by invoking resolve
func (r *Reference) ResolveExample(resolve func(ref string) (*ExampleObj, error)) (*ExampleObj, error) {
	return resolve(r.Ref)
}

// ResolveHeader resolves r by invoking resolve
func (r *Reference) ResolveHeader(resolve func(ref string) (*HeaderObj, error)) (*HeaderObj, error) {
	return resolve(r.Ref)
}

// ResolveRequestBody resolves r by invoking resolve
func (r *Reference) ResolveRequestBody(resolve func(ref string) (*RequestBodyObj, error)) (*RequestBodyObj, error) {
	return resolve(r.Ref)
}

// ResolveCallback resolves r by invoking resolve
func (r *Reference) ResolveCallback(resolve func(ref string) (*CallbackObj, error)) (*CallbackObj, error) {
	return resolve(r.Ref)
}

// ResolvePath resolves r by invoking resolve
func (r *Reference) ResolvePath(resolve func(ref string) (*PathObj, error)) (*PathObj, error) {
	return resolve(r.Ref)
}

// ResolveSecurityScheme resolves r by invoking resolve
func (r *Reference) ResolveSecurityScheme(resolve func(ref string) (*SecuritySchemeObj, error)) (*SecuritySchemeObj, error) {
	return resolve(r.Ref)
}

// ResolveLink resolves r by invoking resolve
func (r *Reference) ResolveLink(resolve func(ref string) (*LinkObj, error)) (*LinkObj, error) {
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
