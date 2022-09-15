package openapi

import "github.com/chanced/uri"

type Ref interface {
	RefURI() *uri.URI
	IsResolved() bool
	Kind() Kind
	RefDst() []interface{}
}

var (
	_ Ref = (*SchemaRef)(nil)
	_ Ref = (*Component[*Server])(nil)
	_ Ref = (*OperationRef)(nil)
)
