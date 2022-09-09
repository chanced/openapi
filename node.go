package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/chanced/jay"
	"github.com/chanced/why"
	"gopkg.in/yaml.v3"
)

type Node interface {
	// MarshalJSON marshals JSON
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals JSON
	UnmarshalJSON(data []byte) error

	Kind() Kind
}

// Validate attempts to validate data which should be an OpenAPI 3.1 definition
// in either json or yaml.
func Validate(data []byte) error {
	var spec map[string]any
	if jay.IsObject(data) {
		err := json.Unmarshal(data, &spec)
		if err == nil {
			return schemas.openapi3_1.OpenAPI.Validate(spec)
		}
	}

	var yn interface{}
	if err := yaml.Unmarshal(data, &yn); err != nil {
		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
	}

	y, err := yaml.Marshal(yn)
	if err != nil {
		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
	}
	j, err := why.YAMLToJSON(y)
	if err != nil {
		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
	}
	err = json.Unmarshal(j, &spec)
	if err != nil {
		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
	}
	return schemas.openapi3_1.OpenAPI.Validate(spec)
}
