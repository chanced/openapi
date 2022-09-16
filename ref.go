package openapi

import "github.com/chanced/uri"

type Ref interface {
	URI() *uri.URI
	IsResolved() bool
	Kind() Kind
	Resolved() Node
}

type ref interface {
	Ref
	resolve(v Node) error
}

var (
	_ Ref = (*SchemaRef)(nil)
	_ ref = (*SchemaRef)(nil)

	_ Ref = (*Reference)(nil)
	_ ref = (*Reference)(nil)

	_ Ref = (*OperationRef)(nil)
	_ ref = (*OperationRef)(nil)
)
