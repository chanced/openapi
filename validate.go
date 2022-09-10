package openapi

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/chanced/jay"
	"github.com/chanced/why"
	"gopkg.in/yaml.v3"
)

// Validate attempts to validate data which should be an OpenAPI 3.1 document
// in either json or yaml.
//
// Validation errors will be
func Validate(data []byte) error {
	var spec map[string]any
	if jay.IsObject(data) {
		err := json.Unmarshal(data, &spec)
		if err == nil {
			return schemas.openapi31[kindDocument].Validate(spec)
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
	return schemas.openapi31[kindDocument].Validate(spec)
}

func validateComponent(kind kind, data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	schema, ok := schemas.openapi31[kind]
	if !ok {
		return errors.New("validation error: unknown component kind")
	}
	return schema.Validate(v)
}
