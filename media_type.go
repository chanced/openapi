package openapi

import (
	"encoding/json"

	"github.com/chanced/openapi/yamlutil"
	"github.com/tidwall/gjson"
)

// Content is a map containing descriptions of potential response payloads. The key is
// a media type or media type range and the value describes it. For
// responses that match multiple keys, only the most specific key is
// applicable. e.g. text/plain overrides text/*
type Content map[string]*MediaType

// MediaType  provides schema and examples for the media type identified by its key.
type MediaType struct {
	//  The schema defining the content of the request, response, or parameter.
	Schema *SchemaObj `json:"schema,omitempty"`
	// Example of the media type. The example object SHOULD be in the correct
	// format as specified by the media type. The example field is mutually
	// exclusive of the examples field. Furthermore, if referencing a schema
	// which contains an example, the example value SHALL override the example
	// provided by the schema.
	Example json.RawMessage `json:"example,omitempty"`
	// Examples of the media type. Each example object SHOULD match the media
	// type and specified schema if present. The examples field is mutually
	// exclusive of the example field. Furthermore, if referencing a schema
	// which contains an example, the examples value SHALL override the example
	// provided by the schema.
	Examples Examples `json:"examples,omitempty"`
	// A map between a property name and its encoding information. The key,
	// being the property name, MUST exist in the schema as a property. The
	// encoding object SHALL only apply to requestBody objects when the media
	// type is multipart or application/x-www-form-urlencoded.
	Encoding   Encodings `json:"encoding,omitempty"`
	Extensions `json:"-"`
}

// MarshalJSON marshals mt into JSON
func (mt MediaType) MarshalJSON() ([]byte, error) {
	type mediatype MediaType
	return marshalExtendedJSON(mediatype(mt))
}

// UnmarshalJSON unmarshals json into mt
func (mt *MediaType) UnmarshalJSON(data []byte) error {
	type mediatype struct {
		Schema     *SchemaObj      `json:"-"`
		Example    json.RawMessage `json:"example,omitempty"`
		Examples   Examples        `json:"examples,omitempty"`
		Encoding   Encodings       `json:"encoding,omitempty"`
		Extensions `json:"-"`
	}

	v := mediatype{}
	if err := unmarshalExtendedJSON(data, &v); err != nil {
		return err
	}
	*mt = MediaType(v)

	g := gjson.GetBytes(data, "schema")
	if g.Exists() {
		s, err := unmarshalSchemaJSON([]byte(g.Raw))
		if err != nil {
			return err
		}
		mt.Schema = s
	}
	return nil
}

// MarshalYAML first marshals and unmarshals into JSON and then marshals into
// YAML
func (mt MediaType) MarshalYAML() (interface{}, error) {
	b, err := json.Marshal(mt)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	return v, err
}

// UnmarshalYAML unmarshals yaml into mt
func (mt *MediaType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return yamlutil.Unmarshal(unmarshal, mt)
}
