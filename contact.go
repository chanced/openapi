package openapi

import "github.com/chanced/uri"

// Contact information for the exposed API.
type Contact struct {
	// The identifying name of the contact person/organization.
	Name Text `json:"name,omitempty"`
	// The URL pointing to the contact information. This MUST be in the form of
	// a URL.
	URL *uri.URI `json:"url,omitempty"`
	// The email address of the contact person/organization. This MUST be in the
	// form of an email address.
	Emails     Text `json:"email,omitempty"`
	Extensions `json:"-"`
}
type contact Contact

// MarshalJSON marshals JSON
func (c Contact) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(contact(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Contact) UnmarshalJSON(data []byte) error {
	var v contact
	err := unmarshalExtendedJSON(data, &v)
	*c = Contact(v)
	return err
}
