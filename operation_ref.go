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
