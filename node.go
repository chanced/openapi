package openapi

import (
	"github.com/chanced/jsonpointer"
	"github.com/chanced/uri"
	"gopkg.in/yaml.v3"
)

type refable interface {
	node
	refable()
}

type Node interface {
	// AbsoluteLocation returns the absolute path of the node in URI form.
	// This includes the URI path of the resource and the JSON pointer
	// of the node.
	//
	// e.g. openapi.json#/components/schemas/Example
	AbsoluteLocation() uri.URI

	// RelativeLocation returns the path as a JSON pointer for the Node.
	RelativeLocation() jsonpointer.Pointer

	// Kind returns the Kind for the given Node
	Kind() Kind

	// Anchors returns a list of all Anchors in the Node and all descendants.
	Anchors() (*Anchors, error)

	// Refs returns a list of all Refs from the Node and all descendants.
	Refs() []Ref

	// MarshalJSON marshals JSON
	//
	// MarshalJSON satisfies the json.Marshaler interface
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	//
	// UnmarshalJSON satisfies the json.Unmarshaler interface
	UnmarshalJSON(data []byte) error

	// UnmarshalYAML satisfies gopkg.in/yaml.v3 Marshaler interface
	MarshalYAML() (interface{}, error)

	// UnmarshalYAML satisfies gopkg.in/yaml.v3 Unmarshaler interface
	UnmarshalYAML(value *yaml.Node) error

	// ResolveNodeByPointer resolves a Node by a jsonpointer. It validates the
	// pointer and then attempts to resolve the Node.
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
	// ResolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error)
}

type node interface {
	Node
	setLocation(loc Location) error
	// init(ctx context.Context, resolver *resolver) error
	// resolveNodeByPointer(ctx context.Context, resolver *resolver, p jsonpointer.Pointer) (node, error)
	mapKind() Kind
	sliceKind() Kind

	// resolveNodeByPointer(ptr jsonpointer.Pointer) (Node, error)
	location() Location
	isNil() bool
	nodes() []node
}

type objSlicedNode interface {
	node
	objSliceKind() Kind
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
