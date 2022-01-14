package openapi

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
