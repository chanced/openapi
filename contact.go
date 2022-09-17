package openapi

import (
	"encoding/json"

	"github.com/chanced/transcodefmt"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

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

func (c Contact) MarshalYAML() (interface{}, error) {
	j, err := c.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcodefmt.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Contact) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcodefmt.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, c)
}
