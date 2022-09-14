package openapi

import "github.com/chanced/uri"

type Ref interface {
	RefURI() *uri.URI
	IsResolved() bool
	IsReference() bool
}

// var (
// 	_ Ref = (*SchemaRef)(nil)
// 	_ Ref = (*Component[*Server])(nil)
// )
