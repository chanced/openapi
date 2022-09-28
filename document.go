package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

// Document root object of the Document document.
type Document struct {
	Location   `json:"-"`
	Extensions `json:"-"`

	// OpenAPI - The OpenAPI Version
	//
	// This string MUST be the version number of the OpenAPI
	// Specification that the OpenAPI document uses. The openapi field SHOULD be
	// used by tooling to interpret the OpenAPI document. This is not related to
	// the API info.version string.
	//
	// 	*required*
	OpenAPI *semver.Version `json:"openapi"`

	// Provides metadata about the API. The metadata MAY be used by
	// tooling as required.
	//
	// 	*required*
	Info *Info `json:"info"`

	// The default value for the $schema keyword within Schema Objects contained
	// within this OAS document.
	JSONSchemaDialect *uri.URI `json:"jsonSchemaDialect,omitempty"`

	// A list of tags used by the document with additional metadata. The order
	// of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared. The
	// tags that are not declared MAY be organized randomly or based on the
	// tools’ logic. Each tag name in the list MUST be unique.
	Tags *TagSlice `json:"tags,omitempty"`

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
	Security *SecurityRequirementSlice `json:"security,omitempty"`

	// Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

func (*Document) Kind() Kind { return KindDocument }

func (d *Document) Refs() []Ref {
	var refs []Ref
	if d.Info != nil {
		refs = append(refs, d.Info.Refs()...)
	}
	if d.Tags != nil {
		refs = append(refs, d.Tags.Refs()...)
	}
	if d.Servers != nil {
		refs = append(refs, d.Servers.Refs()...)
	}
	if d.Paths != nil {
		refs = append(refs, d.Paths.Refs()...)
	}
	if d.Webhooks != nil {
		refs = append(refs, d.Webhooks.Refs()...)
	}
	if d.Components != nil {
		refs = append(refs, d.Components.Refs()...)
	}
	if d.ExternalDocs != nil {
		refs = append(refs, d.ExternalDocs.Refs()...)
	}
	if d.Security != nil {
		refs = append(refs, d.Security.Refs()...)
	}
	return refs
}

func (d *Document) nodes() []node {
	edges := appendEdges(nil, d.Info)
	edges = appendEdges(edges, d.Tags)
	edges = appendEdges(edges, d.Servers)
	edges = appendEdges(edges, d.Paths)
	edges = appendEdges(edges, d.Webhooks)
	edges = appendEdges(edges, d.Components)
	edges = appendEdges(edges, d.Security)
	edges = appendEdges(edges, d.ExternalDocs)
	return edges
}

func (d *Document) isNil() bool {
	return d == nil
}

func (*Document) mapKind() Kind { return KindUndefined }

func (d *Document) setLocation(loc Location) error {
	if d == nil {
		return fmt.Errorf("cannot set location on nil Document")
	}
	d.Location = loc

	if err := d.Info.setLocation(loc.AppendLocation("info")); err != nil {
		return err
	}
	if err := d.Tags.setLocation(loc.AppendLocation("tags")); err != nil {
		return err
	}
	if err := d.Servers.setLocation(loc.AppendLocation("servers")); err != nil {
		return err
	}
	if err := d.Paths.setLocation(loc.AppendLocation("paths")); err != nil {
		return err
	}
	if err := d.Webhooks.setLocation(loc.AppendLocation("webhooks")); err != nil {
		return err
	}
	if err := d.Components.setLocation(loc.AppendLocation("components")); err != nil {
		return err
	}
	if err := d.Security.setLocation(loc.AppendLocation("security")); err != nil {
		return err
	}
	if err := d.ExternalDocs.setLocation(loc.AppendLocation("externalDocs")); err != nil {
		return err
	}
	return nil
}

func (*Document) sliceKind() Kind { return KindUndefined }

// MarshalJSON marshals JSON
func (d Document) MarshalJSON() ([]byte, error) {
	type document Document
	return marshalExtendedJSON(document(d))
}

// UnmarshalJSON unmarshals JSON
func (d *Document) UnmarshalJSON(data []byte) error {
	type openapi Document
	v := openapi{}
	err := unmarshalExtendedJSON(data, &v)
	*d = Document(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (d Document) MarshalYAML() (interface{}, error) {
	j, err := d.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (d *Document) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, d)
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
	if anchors, err = anchors.merge(d.Servers.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.Tags.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.Security.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.Info.Anchors()); err != nil {
		return nil, err
	}
	if anchors, err = anchors.merge(d.ExternalDocs.Anchors()); err != nil {
		return nil, err
	}
	return anchors, nil
}

var _ node = (*Document)(nil)

// func (d *Document) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return d.resolveNodeByPointer(ptr)
// }

// func (d *Document) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return d, nil
// 	}
// 	nxt, tok, _ := ptr.Next()
// 	switch tok {
// 	case "tags":
// 		if d.Tags == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Tags.resolveNodeByPointer(nxt)
// 	case "servers":
// 		if d.Servers == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Servers.resolveNodeByPointer(nxt)
// 	case "paths":
// 		if d.Paths == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Paths.resolveNodeByPointer(nxt)
// 	case "webhooks":
// 		if d.Webhooks == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Webhooks.resolveNodeByPointer(nxt)
// 	case "components":
// 		if d.Components == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Components.resolveNodeByPointer(nxt)
// 	case "externalDocs":
// 		if d.ExternalDocs == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.ExternalDocs.resolveNodeByPointer(nxt)
// 	case "info":
// 		if d.Info == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Info.resolveNodeByPointer(nxt)
// 	case "security":
// 		if d.Security == nil {
// 			return nil, newErrNotFound(d.AbsoluteLocation(), tok)
// 		}
// 		return d.Security.resolveNodeByPointer(nxt)
// 	default:
// 		return nil, newErrNotResolvable(d.AbsoluteLocation(), tok)
// 	}
// }
