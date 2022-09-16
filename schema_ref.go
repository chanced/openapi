package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
	"github.com/chanced/uri"
)

type SchemaRef struct {
	Location
	Ref    *uri.URI `json:"-"`
	Schema *Schema  `json:"-"`
}

func (sr *SchemaRef) Edges() []Node {
	if sr == nil {
		return nil
	}
	return downcastNodes(sr.edges())
}
func (sr *SchemaRef) edges() []node { return []node{sr.Schema} }

func (*SchemaRef) Refs() []Ref { return nil }

func (sr *SchemaRef) IsResolved() bool {
	return sr.Schema != nil
}

func (sr *SchemaRef) URI() *uri.URI { return sr.Ref }

func (*SchemaRef) Kind() Kind      { return KindSchemaRef }
func (*SchemaRef) mapKind() Kind   { return KindUndefined }
func (*SchemaRef) sliceKind() Kind { return KindUndefined }

func (sr *SchemaRef) Resolved() Node {
	if sr == nil {
		return nil
	}
	return sr.Schema
}

func (sr *SchemaRef) resolve(n Node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}

	if s, ok := n.(*Schema); ok {
		sr.Schema = s
		return nil
	}
	return fmt.Errorf("openapi: cannot resolve %s to SchemaRef", n.Kind())
}

func (*SchemaRef) Anchors() (*Anchors, error) { return nil, nil }

func (sr *SchemaRef) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return sr.resolveNodeByPointer(ptr)
}

func (sr *SchemaRef) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	tok, _ := ptr.NextToken()
	if !ptr.IsRoot() {
		if sr.Ref != nil {
			return nil, newErrNotResolvable(sr.Location.AbsolutePath(), tok)
		}
	}
	return sr, nil
}

func (sr *SchemaRef) setLocation(l Location) error {
	if sr == nil {
		return nil
	}
	if sr.Schema != nil {
		if sr.Ref != nil {
			nl, err := NewLocation(*sr.Ref)
			if err != nil {
				return err
			}
			sr.Schema.setLocation(nl)
			return nil
		}
		return sr.Schema.setLocation(l)
	}
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
	sr.Schema = &s
	return err
}

func (sr SchemaRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(sr.Ref)
}

func (sr *SchemaRef) isNil() bool { return sr == nil }

var (
	_ node = (*SchemaRef)(nil)
	// _ Walker = (*SchemaRef)(nil)
	_ Ref = (*SchemaRef)(nil)
)
