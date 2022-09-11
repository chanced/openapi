package openapi

type Node interface {
	Kind() Kind
}
type node interface {
	Node
	// MarshalJSON marshals JSON
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	UnmarshalJSON(data []byte) error

	setLocation(loc Location) error
	// init(ctx context.Context, resolver *resolver) error
	// resolve(ctx context.Context, resolver *resolver, p jsonpointer.Pointer) (node, error)
}
