package openapi

// License information for the exposed API.
type License struct {
	// The license name used for the API.
	//
	// 	*required*
	Name string `json:"name" yaml:"name"`

	// An SPDX license expression for the API. The identifier field is mutually
	// exclusive of the url field.
	Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	// A URL to the license used for the API. This MUST be in the form of a URL.
	// The url field is mutually exclusive of the identifier field.
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
}

func (*License) Nodes() Nodes { return nil }

// Kind returns KindLicense
func (*License) Kind() Kind {
	return KindLicense
}

var _ Node = (*License)(nil)
