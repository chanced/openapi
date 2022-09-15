package openapi

import (
	"encoding/json"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
)

type OperationRef struct {
	Location
	Ref       *uri.URI
	Operation *Operation
}

func (or *OperationRef) Edges() []Node {
	if or == nil {
		return nil
	}
	return downcastNodes(or.edges())
}

func (or *OperationRef) edges() []node {
	if or == nil {
		return nil
	}
	var edges []node
	return appendEdges(edges, or.Operation)
}

func (or *OperationRef) IsRef() bool { return true }

func (or *OperationRef) refs() []node {
	return []node{or.Operation}
}

func (or *OperationRef) Refs() []Ref {
	return nil
}

func (or *OperationRef) IsResolved() bool {
	return or.Operation != nil
}

// RefDst implements Ref
func (or *OperationRef) RefDst() []any {
	return []any{&or.Operation}
}

// RefURI implements Ref
func (or *OperationRef) RefURI() *uri.URI {
	return or.Ref
}

// func (or *OperationRef) Walk(v Visitor) error {
// 	if v == nil {
// 		return nil
// 	}
// 	v, err := v.Visit(or)
// 	if err != nil {
// 		return err
// 	}
// 	if v == nil {
// 		return nil
// 	}
// 	_, err = v.VisitOperationRef(or)
// 	return err
// }

// Anchors implements node
func (*OperationRef) Anchors() (*Anchors, error) { return nil, nil }

func (*OperationRef) Kind() Kind      { return KindOperationRef }
func (*OperationRef) mapKind() Kind   { return KindUndefined }
func (*OperationRef) sliceKind() Kind { return KindUndefined }

func (or *OperationRef) ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if err := ptr.Validate(); err != nil {
		return nil, err
	}
	return or.resolveNodeByPointer(ptr)
}

func (or *OperationRef) resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error) {
	if ptr.IsRoot() {
		return or, nil
	}
	tok, _ := ptr.NextToken()
	return nil, newErrNotResolvable(or.AbsoluteLocation(), tok)
}

func (or OperationRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(or.Ref)
}

func (or *OperationRef) UnmarshalJSON(data []byte) error {
	var uri uri.URI
	if err := json.Unmarshal(data, &uri); err != nil {
		return err
	}
	or.Ref = &uri
	return nil
}

func (o *OperationRef) isNil() bool { return o == nil }
func (op *OperationRef) setLocation(loc Location) error {
	op.Location = loc
	return nil
}

var (
	_ node   = (*OperationRef)(nil)
	_ Walker = (*OperationRef)(nil)
)
