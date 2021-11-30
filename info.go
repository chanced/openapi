package openapi

import "github.com/chanced/openapi/yamlutil"

// Info provides metadata about the API. The metadata MAY be used by the clients
// if needed, and MAY be presented in editing or documentation generation tools
// for convenience.
type Info struct {
	// The title of the API.
	//
	// 	*required*
	Title string `json:"title" yaml:"title"`
	// A short summary of the API.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A description of the API. CommonMark syntax MAY be used for rich text
	// representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A URL to the Terms of Service for the API. This MUST be in the form of a URL.
	TermsOfService string `json:"termsOfService,omitempty" bson:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" bson:"contact,omitempty"`
	// License information for the exposed API.
	License *License `json:"license,omitempty" bson:"license,omitempty"`
	// Version of the OpenAPI document (which is distinct from the OpenAPI
	// Specification version or the API implementation version).
	//
	// 	*required*
	Version    string `json:"version" yaml:"version"`
	Extensions `json:"-"`
}

type info Info

// MarshalJSON marshals JSON
func (i Info) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(info(i))
}

// UnmarshalJSON unmarshals JSON
func (i *Info) UnmarshalJSON(data []byte) error {
	var v info
	err := unmarshalExtendedJSON(data, &v)
	*i = Info(v)
	return err
}

// MarshalYAML marshals YAML
func (i Info) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(i)
}

// UnmarshalYAML unmarshals YAML data into i
func (i *Info) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, i)
}
