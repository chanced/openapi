package openapi

import (
	"bytes"

	"github.com/chanced/jsonx"
	"github.com/chanced/uri"
)

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
	URL *uri.URI `json:"url,omitempty"`
}

// MarshalJSON marshals JSON
func (l License) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	if l.Name != "" {
		b.WriteString(`"name":`)
		jsonx.EncodeAndWriteString(&b, string(l.Name))
	}
	if l.Identifier != "" {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"identifier":`)
		jsonx.EncodeAndWriteString(&b, string(l.Identifier))
	}
	if l.URL != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"url":`)
		bb, err := l.URL.MarshalText()
		if err != nil {
			return nil, err
		}
		jsonx.EncodeAndWriteString(&b, string(bb))
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
