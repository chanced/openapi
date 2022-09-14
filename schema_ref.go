package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/jsonx"
	"github.com/chanced/uri"
)

type SchemaRef struct {
	Location
	Ref      *uri.URI `json:"-"`
	Resolved *Schema  `json:"-"`
}

func (sr *SchemaRef) URI() *uri.URI { return sr.Ref }

func (*SchemaRef) Kind() Kind      { return KindSchemaRef }
func (*SchemaRef) mapKind() Kind   { return KindUndefined }
func (*SchemaRef) sliceKind() Kind { return KindUndefined }

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
			return nil, newErrNotResolvable(sr.Location.AbsoluteLocation(), tok)
		}
	}
	return sr, nil
}

func (sr *SchemaRef) setLocation(l Location) error {
	if sr == nil {
		return nil
	}
	if sr.Resolved != nil {
		if sr.Ref != nil {
			nl, err := NewLocation(*sr.Ref)
			if err != nil {
				return err
			}
			sr.Resolved.setLocation(nl)
			return nil
		}
		return sr.Resolved.setLocation(l)
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
	sr.Resolved = &s
	return err
}

func (sr SchemaRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(sr.Ref)
}

func (sr *SchemaRef) isNil() bool { return sr == nil }

var _ node = (*SchemaRef)(nil)
