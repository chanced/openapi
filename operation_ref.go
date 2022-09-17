package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/jsonpointer"
	"github.com/chanced/transcodefmt"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
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

func (or *OperationRef) Resolved() Node {
	return or.Operation
}

func (or *OperationRef) edges() []node {
	if or == nil {
		return nil
	}
	var edges []node
	return appendEdges(edges, or.Operation)
}

func (or *OperationRef) refs() []node {
	return []node{or.Operation}
}

func (or *OperationRef) Refs() []Ref {
	return nil
}

func (or *OperationRef) IsResolved() bool {
	return or.Operation != nil
}

// URI returns the reference URI
func (or *OperationRef) URI() *uri.URI {
	return or.Ref
}

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
	return nil, newErrNotResolvable(or.AbsolutePath(), tok)
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

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
func (or OperationRef) MarshalYAML() (interface{}, error) {
	j, err := or.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return transcodefmt.YAMLFromJSON(j)
}

// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
func (or *OperationRef) UnmarshalYAML(value *yaml.Node) error {
	j, err := transcodefmt.YAMLFromJSON([]byte(value.Value))
	if err != nil {
		return err
	}
	return json.Unmarshal(j, or)
}

func (o *OperationRef) isNil() bool { return o == nil }
func (op *OperationRef) setLocation(loc Location) error {
	op.Location = loc
	return nil
}

func (o *OperationRef) resolve(n Node) error {
	if o == nil {
		return fmt.Errorf("openapi: OperationRef is nil")
	}
	if n == nil {
		return fmt.Errorf("openapi: node is nil")
	}

	switch n.Kind() {
	case KindOperation:
		o.Operation = n.(*Operation)
	default:
		return fmt.Errorf("openapi: cannot resolve %s to %s", n.Kind(), o.Kind())
	}

	if op, ok := n.(*Operation); ok {
		o.Operation = op
		return nil
	}

	return fmt.Errorf("openapi: failed convert %s to %s", n.Kind(), o.Kind())
}

var (
	_ node = (*OperationRef)(nil)
	_ Ref  = (*OperationRef)(nil)
	_ ref  = (*OperationRef)(nil)
)

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
