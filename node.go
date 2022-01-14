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
