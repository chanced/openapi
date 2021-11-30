package openapi

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/chanced/openapi/yamlutil"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

//go:embed schema
var schemaDir embed.FS
var schemas = map[string]io.ReadCloser{}
var jsonschemaval *jsonschema.Schema
var _ = func() error {

	if err := loadSchemas(); err != nil {
		log.Fatal(err)
	}
	jsonschema.Loaders["https"] = loadSchema
	sch, err := jsonschema.Compile("https://spec.openapis.org/oas/3.1/schema/2021-09-28")
	if err != nil {
		log.Fatal(err)
	}
	jsonschemaval = sch
	return nil
}()

func loadSchema(url string) (io.ReadCloser, error) {
	s, ok := schemas[url]
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", url)
	}
	return s, nil

}

func loadSchemas() error {
	return fs.WalkDir(schemaDir, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || filepath.Ext(d.Name()) != ".json" {
			return nil
		}
		f, err := schemaDir.Open(path)
		if err != nil {
			return err
		}
		j, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		g := gjson.GetBytes(j, "$id")
		if !g.Exists() || len(g.String()) == 0 {
			return errors.New("schema is missing $id")
		}
		schemas[g.String()] = ioutil.NopCloser(bytes.NewReader(j))
		return nil
	})
}

// Validate unmarshals and validates either a single OpenAPI 3.1 specification
// or an array of OpenAPI 3.1 specifications.
//
// The input data can either be a single OpenAPI specification or an array.
func Validate(data []byte) error {
	var list []map[string]interface{}
	d := bytes.TrimSpace(data)
	if len(d) == 0 {
		return errors.New("data may not be empty")
	}
	switch d[0] {
	case '{':
		var o map[string]interface{}
		if err := json.Unmarshal(data, &o); err != nil {
			return err
		}
		list = []map[string]interface{}{o}
	case '[':
		if err := json.Unmarshal(data, &list); err != nil {
			return err
		}
	default:
		b, err := yamlutil.JSONToYAML(data)
		if err != nil {
			return err
		}
		return Validate(b)
	}
	for _, o := range list {
		err := jsonschemaval.Validate(o)
		if err != nil {
			return err
		}
	}
	return nil
}
