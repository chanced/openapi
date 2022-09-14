package openapi

import (
	"bytes"

	"github.com/chanced/jsonx"
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
	// A URL to the Terms of Service for the API. This MUST be in the form of a URL.
	TermsOfService Text `json:"termsOfService,omitempty" bson:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" bson:"contact,omitempty"`
	// License information for the exposed API.
	License *License `json:"license,omitempty" bson:"license,omitempty"`

	Extensions `json:"-"`
}

type info Info

// MarshalJSON marshals JSON
func (i Info) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	if i.Version != "" {
		b.WriteString(`"version":`)
		jsonx.EncodeAndWriteString(&b, string(i.Version))
	}
	if i.Title != "" {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"title":`)
		jsonx.EncodeAndWriteString(&b, string(i.Title))
	}

	if i.Summary != "" {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"summary":`)
		jsonx.EncodeAndWriteString(&b, string(i.Summary))
	}
	if i.Description != "" {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"description":`)
		jsonx.EncodeAndWriteString(&b, string(i.Description))
	}
	if i.TermsOfService != "" {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"termsOfService":`)
		jsonx.EncodeAndWriteString(&b, string(i.TermsOfService))
	}
	if i.Contact != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"contact":`)
		bb, err := i.Contact.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(bb)
	}
	if i.License != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"license":`)
		bb, err := i.License.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(bb)
	}
	return marshalExtensionsInto(&b, i.Extensions)
}

// UnmarshalJSON unmarshals JSON
func (i *Info) UnmarshalJSON(data []byte) error {
	var v info
	err := unmarshalExtendedJSON(data, &v)
	*i = Info(v)
	return err
}
