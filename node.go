package openapi

import (
	"context"

	"github.com/chanced/jsonpointer"
)

type node interface {
	// MarshalJSON marshals JSON
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	UnmarshalJSON(data []byte) error

	kind() kind
	setLocation(loc Location) error
	init(ctx context.Context, resolver *resolver) error
	resolve(ctx context.Context, resolver *resolver, p jsonpointer.Pointer) (node, error)
}
