package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/jsonx"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

type SchemaRefType uint8

const (
	SchamRefTypeUndefined SchemaRefType = iota
	SchemaRefTypeRef
	SchemaRefTypeDynamic
	SchemaRefTypeRecursive
)

type SchemaRef struct {
	Location
	Ref      *uri.URI `json:"-"`
	Resolved *Schema  `json:"-"`

	SchemaRefKind SchemaRefType `json:"-"`
}

func (sr *SchemaRef) Nodes() []Node {
	if sr == nil {
		return nil
	}
	return downcastNodes(sr.nodes())
}

func (sr *SchemaRef) RefType() RefType {
	switch sr.SchemaRefKind {
	case SchemaRefTypeRef:
		return RefTypeSchema
	case SchemaRefTypeDynamic:
		return RefTypeSchemaDynamicRef
	case SchemaRefTypeRecursive:
		return RefTypeSchemaRecursiveRef
	default:
		return RefTypeUndefined
	}
}

func (sr *SchemaRef) RefKind() Kind { return KindSchema }

func (sr *SchemaRef) nodes() []node { return []node{sr.Resolved} }

func (*SchemaRef) Refs() []Ref { return nil }

func (sr *SchemaRef) IsResolved() bool {
	return sr.Resolved != nil
}

func (sr *SchemaRef) URI() *uri.URI { return sr.Ref }

func (*SchemaRef) Kind() Kind      { return KindSchemaRef }
func (*SchemaRef) mapKind() Kind   { return KindUndefined }
func (*SchemaRef) sliceKind() Kind { return KindUndefined }

func (sr *SchemaRef) ResolvedNode() Node {
	if sr == nil {
		return nil
	}
	return sr.Resolved
}

// func (sr *SchemaRef) Clone() *SchemaRef {
// 	if sr == nil {
// 		return nil
// 	}
// 	c := *sr
// 	return &c
// }

func (sr *SchemaRef) resolve(n Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if s, ok := n.(*Schema); ok {
		sr.Resolved = s
		return nil
	}
	return NewResolutionError(sr, KindSchema, n.Kind())
}

func (*SchemaRef) Anchors() (*Anchors, error) { return nil, nil }

func (sr *SchemaRef) setLocation(l Location) error {
	if sr == nil {
		return nil
	}
	sr.Location = l
	// if sr.Schema != nil {
	// 	if sr.Ref != nil {
	// 		nl, err := NewLocation(*sr.Ref)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		sr.Schema.setLocation(nl)
	// 		return nil
	// 	}
	// 	return sr.Schema.setLocation(l)
	// }
	return nil
}

func (sr *SchemaRef) UnmarshalJSON(data []byte) error {
	if jsonx.IsString(data) {
		var u uri.URI
		if err := json.Unmarshal(data, &u); err != nil {
			return err
		}
		sr.Ref = &u
		return nil
	}

	var s Schema
	err := json.Unmarshal(data, &s)
	sr.Resolved = &s
	return err
}

func (sr SchemaRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(sr.Ref)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (sr SchemaRef) MarshalYAML() (interface{}, error) {
	j, err := sr.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcode.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (sr *SchemaRef) UnmarshalYAML(value *yaml.Node) error {
	v, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	j, err := transcode.JSONFromYAML(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, sr)
}

func (sr *SchemaRef) isNil() bool { return sr == nil }

func (sr *SchemaRef) Clone() *SchemaRef {
	if sr == nil {
		return nil
	}
	var ref *uri.URI
	if sr.Ref != nil {
		ref = sr.Ref.Clone()
	}
	return &SchemaRef{
		Ref: ref,
		Location: Location{
			absolute: sr.Location.absolute,
			relative: sr.Location.relative,
		},
		Resolved:      sr.Resolved.Clone(), // should this be cloned?
		SchemaRefKind: sr.SchemaRefKind,
	}
}

var (
	_ node = (*SchemaRef)(nil)

	_ Ref = (*SchemaRef)(nil)
)

// func (sr *SchemaRef) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	if err := ptr.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return sr.resolveNodeByPointer(ptr)
// }

// func (sr *SchemaRef) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
// 	tok, _ := ptr.NextToken()
// 	if !ptr.IsRoot() {
// 		if sr.Ref != nil {
// 			return nil, newErrNotResolvable(sr.Location.AbsoluteLocation(), tok)
// 		}
// 	}
// 	return sr, nil
// }
