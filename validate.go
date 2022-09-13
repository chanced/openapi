package openapi

// 	y, err := yaml.Marshal(yn)
// 	if err != nil {
// 		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
// 	}
// 	j, err := transcodefmt.YAMLToJSON(y)
// 	if err != nil {
// 		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
// 	}
// 	err = json.Unmarshal(j, &spec)
// 	if err != nil {
// 		return fmt.Errorf("error unmarshaling json or yaml: %w", err)
// 	}
// 	return schemas.openapi31[KindDocument].Validate(spec)
// }

// func validateComponent(kind Kind, data []byte) error {
// 	var v interface{}
// 	if err := json.Unmarshal(data, &v); err != nil {
// 		return err
// 	}
// 	schema, ok := schemas.openapi31[kind]
// 	if !ok {
// 		return errors.New("validation error: unknown component kind")
// 	}
// 	return schema.Validate(v)
// }
