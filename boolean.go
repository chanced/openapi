package openapi

// Boolean is a bool value, which can be used as a Schema
type Boolean bool

// SchemaKind returns SchemaKindBool
func (b Boolean) SchemaKind() SchemaKind {
	return SchemaKindBool
}

// IsRef returns false
func (b Boolean) IsRef() bool {
	return false
}
