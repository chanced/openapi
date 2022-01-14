package yamlutil

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/chanced/dynamic"
	"sigs.k8s.io/yaml"
)

// Unmarshal unmarshals in from YAML by marshaling and unmarshaling to json
func Unmarshal(unmarshal func(in interface{}) error, out interface{}) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	v, err := subset(i)
	if err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	jb, err := YAMLToJSON(b)

	if err != nil {
		return err
	}
	return json.Unmarshal(jb, out)
}

// YAMLToJSON converts YAML to JSON
func YAMLToJSON(data []byte) ([]byte, error) {
	// found sigs.k8s.io/yaml which does a great job of converting this over so
	// I'm just going to use it instead of what I had.
	// TODO: Re-implement YAMLToJSON as sigs/yaml doesn't handle big numbers
	return yaml.YAMLToJSON(data)
}

// Marshal returns an interface{} representation of src to be marshaled into
// YAML
func Marshal(src interface{}) (interface{}, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// JSONToYAML converts JSON to YAML.
func JSONToYAML(data []byte) ([]byte, error) {
	return yaml.JSONToYAML(data)
}

// EncodeYAMLToJSON encodes YAML as JSON and returns an io.Reader of the JSON
func EncodeYAMLToJSON(r io.Reader) (io.Reader, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b, err = YAMLToJSON(b)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

func subset(u interface{}) (interface{}, error) {
	switch t := u.(type) {
	case []interface{}:
		return subsetSlice(t)
	case map[interface{}]interface{}:
		return subsetObj(t)
	case map[string]interface{}:
		return subsetMap(t)
	default:
		return t, nil
	}
}

func subsetSlice(t []interface{}) ([]interface{}, error) {
	res := make([]interface{}, len(t))
	for i, v := range t {
		iv, err := subset(v)
		if err != nil {
			return nil, err
		}
		res[i] = iv
	}
	return res, nil
}
func subsetMap(t map[string]interface{}) (map[string]interface{}, error) {
	// this should be okay but going to check it regardless.
	res := make(map[string]interface{}, len(t))
	for k, v := range t {
		cv, err := subset(v)
		if err != nil {
			return nil, err
		}
		res[k] = cv
	}
	return res, nil
}
func subsetObj(t map[interface{}]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(t))
	for k, v := range t {
		ks := new(dynamic.String)
		err := ks.Set(k)
		if err != nil {
			return nil, err
		}
		cv, err := subset(v)
		if err != nil {
			return nil, err
		}
		res[ks.String()] = cv
	}
	return res, nil
}
