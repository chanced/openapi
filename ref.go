package openapi

import "github.com/chanced/uri"

type RefKind uint8

const (
	RefKindUndefined RefKind = iota
	RefKindReference
	RefKindSchemaRef
	RefKindSchemaDynamicRef
	RefKindSchemaRecursiveRef
	RefKindOperationRef
)

type Ref interface {
	Node
	URI() *uri.URI
	IsResolved() bool
	Kind() Kind
	Resolved() Node
}

type ref interface {
	Ref
	resolve(v Node) error
}

// IsRef returns true for the following types:
//   - *Reference
//   - *SchemaRef
//   - *OperationRef
func IsRef(node Node) bool {
	switch node.Kind() {
	case KindReference, KindSchemaRef, KindOperationRef:
		_, ok := node.(Ref)
		if !ok {
			panic("node is not a Ref. This is a bug. Please report it to github.com/chanced/openapi")
		}
		return true
	default:
		return false
	}
}

var (
	_ Ref = (*SchemaRef)(nil)
	_ ref = (*SchemaRef)(nil)

	_ Ref = (*Reference)(nil)
	_ ref = (*Reference)(nil)

	_ Ref = (*OperationRef)(nil)
	_ ref = (*OperationRef)(nil)
)
