package openapi

type Node interface {
	Kind() Kind
	// Nodes() map[string]Node

	// NodeKind(nodepath string) Kind
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

func (v PassthroughVisitor) Visit(node Node, path string) (Visitor, error) {
	return v, nil
}
