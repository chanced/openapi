package openapi

import (
	"encoding/json"

	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

// Contact information for the exposed API.
type Contact struct {
	Extensions `json:"-"`
	Location   `json:"-"`
	// The identifying name of the contact person/organization.
	Name Text `json:"name,omitempty"`
	// The URL pointing to the contact information. This MUST be in the form of
	// a URL.
	URL *uri.URI `json:"url,omitempty"`
	// The email address of the contact person/organization. This MUST be in the
	// form of an email address.
	Emails Text `json:"email,omitempty"`
}

func (*Contact) Anchors() (*Anchors, error) { return nil, nil }

// Kind returns KindContact
func (*Contact) Kind() Kind { return KindContact }

func (*Contact) Refs() []Ref { return nil }

// func (c *Contact) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return c.resolveNodeByPointer(ptr)
// }

// func (c *Contact) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return c, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	return nil, newErrNotResolvable(c.AbsoluteLocation(), tok)
// }

func (*Contact) nodes() []node        { return nil }
func (c *Contact) isNil() bool        { return c == nil }
func (c *Contact) location() Location { return c.Location }

func (c *Contact) setLocation(loc Location) error {
	if c != nil {
		c.Location = loc
	}
	return nil
}

func (*Contact) sliceKind() Kind { return KindUndefined }
func (*Contact) mapKind() Kind   { return KindUndefined }

// MarshalJSON marshals JSON
func (c Contact) MarshalJSON() ([]byte, error) {
	type contact Contact
	return marshalExtendedJSON(contact(c))
}

// UnmarshalJSON unmarshals JSON
func (c *Contact) UnmarshalJSON(data []byte) error {
	type contact Contact
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
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Contact) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcode.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, c)
}

var _ node = (*Contact)(nil)
