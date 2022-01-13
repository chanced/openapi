package openapi

import "github.com/chanced/openapi/yamlutil"

// Contact information for the exposed API.
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty"`
	// The URL pointing to the contact information. This MUST be in the form of
	// a URL.
	URL string `json:"url,omitempty"`
	// The email address of the contact person/organization. This MUST be in the
	// form of an email address.
	Emails     string `json:"email,omitempty"`
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

// MarshalYAML marshals YAML
func (c Contact) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(c)
}

// UnmarshalYAML unmarshals YAML
func (c *Contact) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, c)
}

// Kind returns KindContact
func (*Contact) Kind() Kind {
	return KindContact
}

func (*Contact) Nodes() Nodes {
	return nil
}

var _ Node = (*Contact)(nil)
