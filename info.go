package openapi

import (
	"encoding/json"

	"github.com/chanced/transcodefmt"
	"gopkg.in/yaml.v3"
)

// Info provides metadata about the API. The metadata MAY be used by the clients
// if needed, and MAY be presented in editing or documentation generation tools
// for convenience.
type Info struct {
	// Version of the OpenAPI document (which is distinct from the OpenAPI
	// Specification version or the API implementation version).
	//
	// 	*required*
	Version Text `json:"version"`
	// The title of the API.
	//
	// 	*required*
	Title Text `json:"title"`

	// A short summary of the API.
	Summary Text `json:"summary,omitempty"`

	// A description of the API. CommonMark syntax MAY be used for rich text
	// representation.
	Description Text `json:"description,omitempty"`

	// A URL to the Terms of Service for the API. This MUST be in the form of a
	// URL.
	TermsOfService Text `json:"termsOfService,omitempty" bson:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" bson:"contact,omitempty"`

	// License information for the exposed API.
	License *License `json:"license,omitempty" bson:"license,omitempty"`

	Extensions `json:"-"`
}

// MarshalJSON marshals JSON
func (i Info) MarshalJSON() ([]byte, error) {
	type info Info

	return marshalExtendedJSON(info(i))
}

// UnmarshalJSON unmarshals JSON
func (i *Info) UnmarshalJSON(data []byte) error {
	type info Info
	var v info
	err := unmarshalExtendedJSON(data, &v)
	*i = Info(v)
	return err
}
// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (i Info) MarshalYAML() (interface{}, error) {
	j, err := i.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcodefmt.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (i *Info) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcodefmt.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, i)
}
