package openapi

import "github.com/chanced/jsonpointer"

type Node interface {
	Kind() Kind
	// ResolveNodeByPointer resolves a Node by a jsonpointer. It validates the pointer and then
	// attempts to resolve the Node.
	//
	// # Errors
	//
	// - [ErrNotFound] indicates that the component was not found
	//
	// - [ErrNotResolvable] indicates that the pointer path can not resolve to a
	// Node
	//
	// - [jsonpointer.ErrMalformedEncoding] indicates that the pointer encoding
	// is malformed
	//
	// - [jsonpointer.ErrMalformedStart] indicates that the pointer is not empty
	// and does not start with a slash
	ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error)

	Anchors() (*Anchors, error)

	Refs() []Ref
}

// IsRef returns true if the Node is any of the following:
//   - Reference
//   - SchemaRef
//   - OperationRef
func IsRef(n Node) bool {
	switch n.Kind() {
	case KindReference, KindSchemaRef, KindOperationRef:
		return true
	default:
		return false
	}
}

type node interface {
	Node
	// MarshalJSON marshals JSON
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	UnmarshalJSON(data []byte) error

	setLocation(loc Location) error
	// init(ctx context.Context, resolver *resolver) error
	// resolveNodeByPointer(ctx context.Context, resolver *resolver, p jsonpointer.Pointer) (node, error)
	mapKind() Kind
	sliceKind() Kind

	resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error)
	location() Location
	isNil() bool
	edges() []node
}

func downcastNodes(n []node) []Node {
	nodes := make([]Node, len(n))
	for i, v := range n {
		nodes[i] = v
	}
	return nodes
}

func appendEdges(nodes []node, elems ...node) []node {
	for _, n := range elems {
		if !n.isNil() {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
