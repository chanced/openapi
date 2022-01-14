package openapi

import (
	"encoding"
	"fmt"
	"strconv"
)

type NodeDetail struct {
	Node
	// TargetKind is the Kind regardless of the Node, which could be the
	// TargetKind or a Reference
	TargetKind Kind
}

func newNodeDetail(node Node, kind Kind) NodeDetail {
	return NodeDetail{
		Node:       node,
		TargetKind: kind,
	}
}

type Nodes map[string]NodeDetail

func makeNodes(nl nodes) Nodes {
	ns := make(Nodes, len(nl))
	for k, n := range nl {
		ns.maybeAdd(k, n.Node, n.TargetKind)
	}
	if len(ns) == 0 {
		return nil
	}
	return ns
}

type lengther interface {
	Len() int
}

type nodes map[interface{}]NodeDetail

func (ns *Nodes) maybeAdd(key interface{}, n Node, kind Kind) {
	var k string
	switch kt := key.(type) {
	case string:
		k = kt
	case int:
		k = strconv.Itoa(kt)
	case fmt.Stringer:
		k = kt.String()
	case encoding.TextMarshaler:
		t, err := kt.MarshalText()
		if err != nil {
			return
		}
		k = string(t)
	default:
		panic("key must be a string, int, or implement either fmt.String or encoding.TextMarshaler")
	}

	if n == nil {
		return
	}
	if l, ok := n.(lengther); ok {
		if l.Len() == 0 {
			return
		}
	}
	if ns == nil {
		*ns = Nodes{
			k: newNodeDetail(n, kind),
		}
		return
	}
	(*ns)[k] = newNodeDetail(n, kind)
}

type Node interface {
	Kind() Kind
	Nodes() Nodes
}

type Visitor interface {
	Visit(node Node, path string) (Visitor, error)
}

// InfoVisitor is implemented by visitors with a VisitInfo method
type InfoVisitor interface {
	Visitor
	VisitInfo(node *Info, path string) (Visitor, error)
}

// ServerVisitor is implemented by a visitor with a VisitServer method
type ServerVisitor interface {
	Visitor
	VisitServer(node *Server, path string) (Visitor, error)
}

// PassthroughVisitor is a noop visitor which passes all nodes through
//
// The use case for this is to implement a visitor which only visits a subset
// of the nodes.
// For example, to only visit the Info node, you could implement InfoVisitor such as:
//	type InfoVisitor struct {
// 		openapi.PassthroughVisitor
// 	}
//
//	func (iv InfoVisitor) VisitInfo(node *Info) (Visitor, error) {
//	// do something with node
// 		return iv, nil
// }
type PassthroughVisitor struct{}

// Visit returns v and nil
func (v PassthroughVisitor) Visit(_ Node, _ string) (Visitor, error) {
	return v, nil
}

// WalkOpts are options for Walk
type WalkOpts struct {
	BasePath *string
}

// WalkOpt is an Option for walk
type WalkOpt func(opts *WalkOpts) *WalkOpts

// WalkWithBasePath sets the base path from which to resolve relative paths
var WalkWithBasePath = func(basePath string) WalkOpt {
	return func(opts *WalkOpts) *WalkOpts {
		if opts == nil {
			*opts = WalkOpts{}
		}
		opts.BasePath = &basePath
		return opts
	}
}

// Walk walks the node tree and calls the visitor for each node
func Walk(node Node, opts ...WalkOpt) error {
	options := &WalkOpts{}
	for _, opt := range opts {
		options = opt(options)
	}
	_ = node
	panic("not impl") // TODO: impl
}
