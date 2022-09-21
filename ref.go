package openapi

import "github.com/chanced/uri"

const (
	RefTypeUndefined RefType = iota
	RefTypeComponent
	RefTypeSchema
	RefTypeSchemaDynamicRef
	RefTypeSchemaRecursiveRef
	RefTypeOperationRef
)

type RefType uint8

func (rk RefType) String() string {
	switch rk {
	case RefTypeComponent:
		return "Reference"
	case RefTypeSchema:
		return "SchemaRef"
	case RefTypeSchemaDynamicRef:
		return "SchemaRef"
	case RefTypeSchemaRecursiveRef:
		return "SchemaRef"
	case RefTypeOperationRef:
		return "OperationRef"
	default:
		return "Undefined"
	}
}

type Ref interface {
	Node
	URI() *uri.URI
	IsResolved() bool
	ResolvedNode() Node
	// ReferencedKind returns the Kind for the referenced node
	RefKind() Kind
	// RefType returns the RefType for the reference
	RefType() RefType
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

	_ Ref = (*Reference[*Response])(nil)
	_ ref = (*Reference[*Response])(nil)

	_ Ref = (*OperationRef)(nil)
	_ ref = (*OperationRef)(nil)
)
