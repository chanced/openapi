package openapi

import (
	"encoding/json"

	"github.com/Masterminds/semver"
	"github.com/chanced/transcode"
	"gopkg.in/yaml.v3"
)

// Info provides metadata about the API. The metadata MAY be used by the clients
// if needed, and MAY be presented in editing or documentation generation tools
// for convenience.
type Info struct {
	Extensions `json:"-"`
	Location   `json:"-"`

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

	// A URL to the Terms of Service for the API. This MUST be in the form of a
	// URL.
	TermsOfService Text `json:"termsOfService,omitempty" bson:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" bson:"contact,omitempty"`

	// License information for the exposed API.
	License *License `json:"license,omitempty" bson:"license,omitempty"`
}

func (*Info) Anchors() (*Anchors, error) { return nil, nil }

func (*Info) Kind() Kind { return KindInfo }

func (*Info) Refs() []Ref { return nil }

// func (i *Info) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	err := ptr.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if ptr.IsRoot() {
// 		return i, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	switch tok {
// 	case "contact":
// 		if i.Contact == nil {
// 			return nil, newErrNotFound(i.absolute, tok)
// 		}
// 		return i.Contact, nil
// 	case "license":
// 		return i.License, nil
// 	}
// 	return nil, newErrNotResolvable(i.AbsoluteLocation(), tok)
// }

// func (i *Info) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if ptr.IsRoot() {
// 		return i, nil
// 	}
// 	tok, _ := ptr.NextToken()
// 	switch tok {
// 	case "contact":
// 		if i.Contact == nil {
// 			return nil, newErrNotFound(i.AbsoluteLocation(), tok)
// 		}
// 		return i.Contact, nil
// 	case "license":
// 		if i.License == nil {
// 			return nil, newErrNotFound(i.AbsoluteLocation(), tok)
// 		}
// 		return i.License, nil
// 	default:
// 		return nil, newErrNotResolvable(i.AbsoluteLocation(), tok)
// 	}
// }

func (i *Info) nodes() []node {
	edges := appendEdges(nil, i.Contact)
	edges = appendEdges(edges, i.License)
	return edges
}

func (i *Info) isNil() bool {
	return i == nil
}

func (i *Info) location() Location {
	return i.Location
}

func (i *Info) SemVer() (*semver.Version, error) {
	return semver.NewVersion(i.Version.String())
}

func (*Info) mapKind() Kind { return KindUndefined }

func (i *Info) setLocation(loc Location) error {
	if i == nil {
		return nil
	}
	i.Location = loc
	if err := i.Contact.setLocation(loc.AppendLocation("contact")); err != nil {
		return err
	}
	if err := i.License.setLocation(loc.AppendLocation("license")); err != nil {
		return err
	}
	return nil
}

func (*Info) sliceKind() Kind { return KindUndefined }

// MarshalJSON marshals JSON
func (i Info) MarshalJSON() ([]byte, error) {
	type info Info

	return marshalExtendedJSON(info(i))
}

// UnmarshalJSON unmarshals JSON
func (i *Info) UnmarshalJSON(data []byte) error {
	type info Info
	var v info
	err := unmarshalExtendedJSON(data, &v)
	*i = Info(v)
	return err
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (i Info) MarshalYAML() (interface{}, error) {
	j, err := i.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (i *Info) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, i)
}

var _ node = (*Info)(nil)
