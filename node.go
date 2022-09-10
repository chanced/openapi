package openapi

import "context"

type node interface {
	// MarshalJSON marshals JSON
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	UnmarshalJSON(data []byte) error

	kind() kind

	resolve(ctx context.Context, resolver resolver, p string, kind kind) (interface{}, error)
}
