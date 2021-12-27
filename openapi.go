package openapi

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/chanced/openapi/yamlutil"
)

type openapi OpenAPI

// OpenAPI root object of the OpenAPI document.
type OpenAPI struct {
	// Version - OpenAPI Version
	//
	// This string MUST be the version number of the OpenAPI
	// Specification that the OpenAPI document uses. The openapi field SHOULD be
	// used by tooling to interpret the OpenAPI document. This is not related to
	// the API info.version string.
	Version string `json:"openapi" yaml:"openapi"`
	// Provides metadata about the API. The metadata MAY be used by
	// tooling as required.
	//
	// 	*required*
	Info *Info `json:"info" yaml:"info"`
	// The default value for the $schema keyword within Schema Objects contained
	// within this OAS document. This MUST be in the form of a URI.
	JSONSchemaDialect string `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	// An array of Server Objects, which provide connectivity information to a
	// target server. If the servers property is not provided, or is an empty
	// array, the default value would be a Server Object with a url value of /.
	Servers []*Server `json:"servers,omitempty" yaml:"servers,omitempty,omtiempty"`
	// The available paths and operations for the API.
	Paths *Paths `json:"paths,omitempty" yaml:"paths,omitempty"`
	// The incoming webhooks that MAY be received as part of this API and that
	// the API consumer MAY choose to implement. Closely related to the
	// callbacks feature, this section describes requests initiated other than
	// by an API call, for example by an out of band registration. The key name
	// is a unique string to refer to each webhook, while the (optionally
	// referenced) Path Item Object describes a request that may be initiated by
	// the API provider and the expected responses. An example is available.
	Webhooks *PathItems `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	// An element to hold various schemas for the document.
	Components *Components `json:"components,omitempty" yaml:"components,omitempty"`
	// A list of tags used by the document with additional metadata. The order
	// of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared. The
	// tags that are not declared MAY be organized randomly or based on the
	// tools’ logic. Each tag name in the list MUST be unique.
	Tags []*Tag `json:"tags,omitempty" yaml:"tags,omitempty"`
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
	Security []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	Extensions `json:"-"`
}

// Kind returns KindOpenAPI
func (*OpenAPI) Kind() Kind {
	return KindOpenAPI
}

// Validate validates an OpenAPI 3.1 specification
func (o OpenAPI) Validate() error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	return validate(m)
}

// MarshalJSON marshals JSON
func (o OpenAPI) MarshalJSON() ([]byte, error) {
	return marshalExtendedJSON(openapi(o))
}

// UnmarshalJSON unmarshals JSON
func (o *OpenAPI) UnmarshalJSON(data []byte) error {
	v := openapi{}
	err := unmarshalExtendedJSON(data, &v)
	*o = OpenAPI(v)
	return err
}

// MarshalYAML marshals o into yaml
func (o OpenAPI) MarshalYAML() (interface{}, error) {
	return yamlutil.Marshal(o)
}

// UnmarshalYAML unmarshals YAML data into o
func (o *OpenAPI) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, o)
}

// // Resolve resolves all references and returns a *ResolvedOpenAPI instance or an
// // error
// func (o *OpenAPI) Resolve(resolver Resolver) (*ResolvedOpenAPI, error) {
// 	// r := &ResolvedOpenAPI{}

// 	// TODO: finish this
// 	panic("not implemented")
// }

// EncodeJSON encodes OpenAPI to JSON
func (o *OpenAPI) EncodeJSON() (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(o)
	return buf, err
}

// ResolvedOpenAPI is an OpenAPI 3.1 specification with resolved references
type ResolvedOpenAPI struct {
	// Version - OpenAPI Version
	//
	// This string MUST be the version number of the OpenAPI
	// Specification that the OpenAPI document uses. The openapi field SHOULD be
	// used by tooling to interpret the OpenAPI document. This is not related to
	// the API info.version string.
	Version string `json:"openapi" yaml:"openapi"`
	// Provides metadata about the API. The metadata MAY be used by
	// tooling as required.
	//
	// 	*required*
	Info *Info `json:"info" yaml:"info"`
	// The default value for the $schema keyword within Schema Objects contained
	// within this OAS document. This MUST be in the form of a URI.
	JSONSchemaDialect string `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	// An array of Server Objects, which provide connectivity information to a
	// target server. If the servers property is not provided, or is an empty
	// array, the default value would be a Server Object with a url value of /.
	Servers []*Server `json:"servers,omitempty" yaml:"servers,omitempty,omtiempty"`
	// The available paths and operations for the API.
	Paths *Paths `json:"paths,omitempty" yaml:"paths,omitempty"`
	// The incoming webhooks that MAY be received as part of this API and that
	// the API consumer MAY choose to implement. Closely related to the
	// callbacks feature, this section describes requests initiated other than
	// by an API call, for example by an out of band registration. The key name
	// is a unique string to refer to each webhook, while the (optionally
	// referenced) Path Item Object describes a request that may be initiated by
	// the API provider and the expected responses. An example is available.
	Webhooks *PathItems `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	// An element to hold various schemas for the document.
	Components *Components `json:"components,omitempty" yaml:"components,omitempty"`
	// A list of tags used by the document with additional metadata. The order
	// of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared. The
	// tags that are not declared MAY be organized randomly or based on the
	// tools’ logic. Each tag name in the list MUST be unique.
	Tags []*Tag `json:"tags,omitempty" yaml:"tags,omitempty"`
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
	Security []*SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	// externalDocs	Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Extensions   `json:"-"`
}

// Kind returns KindResolvedOpenAPI
func (*ResolvedOpenAPI) Kind() Kind {
	return KindResolvedOpenAPI
}

var (
	_ Node = (*OpenAPI)(nil)
	_ Node = (*ResolvedOpenAPI)(nil)
)
