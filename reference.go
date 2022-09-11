package openapi

import (
	"errors"

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
	Ref Text `yaml:"$ref" json:"$ref"`
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
	Location *Location `json:"-"`
}

func (r *Reference) setLocation(loc Location) error {
	if r == nil {
		return nil
	}
	r.Location = &loc
	return nil
}

func (r *Reference) Kind() Kind { return KindReference }

func isRefJSON(data []byte) bool {
	r := gjson.GetBytes(data, "$ref")
	return r.Str != ""
}
