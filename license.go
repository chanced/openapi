package openapi

// License information for the exposed API.
type License struct {
	// The license name used for the API.
	//
	// 	*required*
	Name Text `json:"name"`

	// An SPDX license expression for the API. The identifier field is mutually
	// exclusive of the url field.
	Identifier Text `json:"identifier,omitempty"`
	// A URL to the license used for the API. This MUST be in the form of a URL.
	// The url field is mutually exclusive of the identifier field.
	URL Text `json:"url,omitempty"`
}
