package openapi

import (
	"bytes"
	"encoding/json"

	"github.com/Masterminds/semver"
	"github.com/chanced/jsonx"
	"github.com/chanced/uri"
)

// Document root object of the Document document.
type Document struct {
	// OpenAPI - The OpenAPI Version
	//
	// This string MUST be the version number of the OpenAPI
	// Specification that the OpenAPI document uses. The openapi field SHOULD be
	// used by tooling to interpret the OpenAPI document. This is not related to
	// the API info.version string.
	OpenAPI *semver.Version `json:"openapi"`
	// Provides metadata about the API. The metadata MAY be used by
	// tooling as required.
	//
	// 	*required*
	Info *Info `json:"info"`
	// The default value for the $schema keyword within Schema Objects contained
	// within this OAS document. This MUST be in the form of a URI.
	JSONSchemaDialect *uri.URI `json:"jsonSchemaDialect,omitempty"`
	// An array of Server Objects, which provide connectivity information to a
	// target server. If the servers property is not provided, or is an empty
	// array, the default value would be a Server Object with a url value of /.
	Servers *ServerSlice `json:"servers,omitempty" yaml:"servers,omitempty,omtiempty"`
	// The available paths and operations for the API.
	Paths *Paths `json:"paths,omitempty"`
	// The incoming webhooks that MAY be received as part of this API and that
	// the API consumer MAY choose to implement. Closely related to the
	// callbacks feature, this section describes requests initiated other than
	// by an API call, for example by an out of band registration. The key name
	// is a unique string to refer to each webhook, while the (optionally
	// referenced) Path Item Object describes a request that may be initiated by
	// the API provider and the expected responses. An example is available.
	Webhooks *PathItemMap `json:"webhooks,omitempty"`
	// An element to hold various schemas for the document.
	Components *Components `json:"components,omitempty"`
	// A list of tags used by the document with additional metadata. The order
	// of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared. The
	// tags that are not declared MAY be organized randomly or based on the
	// tools’ logic. Each tag name in the list MUST be unique.
	Tags *TagSlice `json:"tags,omitempty"`
	// A declaration of which security mechanisms can be used across the API.
	//
	// The list of values includes alternative security requirement objects that
	// can be used.
	//
	// Only one of the security requirement objects need to be
	// satisfied to authorize a request. Individual operations can override this
	// definition.
	//
	// To make security optional, an empty security requirement ({})
	// can be included in the array.
	//
	Security []SecurityRequirement `json:"security,omitempty"`
	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	Extensions   `json:"-"`
}

// MarshalJSON marshals JSON
func (d Document) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	b.WriteByte('{')
	if d.OpenAPI != nil {
		d.OpenAPI.MarshalJSON()
		b.WriteString(`"openapi":`)
		b.WriteString("\"" + d.OpenAPI.Original() + "\"")
	}
	if d.Info != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"info":`)
		i, err := d.Info.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(i)
	}
	if d.JSONSchemaDialect != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}

		b.WriteString(`"jsonSchemaDialect":`)
		jsonx.EncodeAndWriteString(&b, d.JSONSchemaDialect.String())
	}

	if d.Servers != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"servers":`)
		s, err := d.Servers.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(s)
	}
	if d.Paths != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"paths":`)
		p, err := d.Paths.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(p)
	}
	if d.Webhooks != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"webhooks":`)
		w, err := d.Webhooks.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(w)
	}

	if d.Security != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"security":`)
		s, err := json.Marshal(d.Security)
		if err != nil {
			return nil, err
		}
		b.Write(s)
	}
	if d.Tags != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"tags":`)
		t, err := d.Tags.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(t)
	}
	if d.ExternalDocs != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"externalDocs":`)
		e, err := d.ExternalDocs.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(e)
	}
	if d.Components != nil {
		if b.Len() > 1 {
			b.WriteByte(',')
		}
		b.WriteString(`"components":`)
		c, err := d.Components.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(c)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// UnmarshalJSON unmarshals JSON
func (d *Document) UnmarshalJSON(data []byte) error {
	type openapi Document
	v := openapi{}
	err := unmarshalExtendedJSON(data, &v)
	*d = Document(v)
	return err
}

func (d *Document) Anchors() (*Anchors, error) {
	if d == nil {
		return nil, nil
	}
	var anchors *Anchors
	var err error
	if anchors, err = anchors.merge(d.Paths.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.Components.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.Webhooks.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}
