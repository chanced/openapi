package openapi

import (
	"embed"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed schema
var schemaDir embed.FS
var compiler = jsonschema.NewCompiler()

var schemas struct {
	schema202012 *jsonschema.Schema
	v3_1         struct {
		OpenAPI      *jsonschema.Schema
		Operation    *jsonschema.Schema
		Callback     *jsonschema.Schema
		Example      *jsonschema.Schema
		Header       *jsonschema.Schema
		License      *jsonschema.Schema
		Link         *jsonschema.Schema
		Parameter    *jsonschema.Schema
		Request      *jsonschema.Schema
		Response     *jsonschema.Schema
		Security     *jsonschema.Schema
		Tag          *jsonschema.Schema
		Path         *jsonschema.Schema
		MediaType    *jsonschema.Schema
		Info         *jsonschema.Schema
		Contact      *jsonschema.Schema
		Schema       *jsonschema.Schema
		XML          *jsonschema.Schema
		Encoding     *jsonschema.Schema
		Reference    *jsonschema.Schema
		ExternalDocs *jsonschema.Schema
	}
}

func init() {
	log.SetFlags(0)
	compiler.Draft = jsonschema.Draft2020
	err := fs.WalkDir(schemaDir, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".json" {
			return err
		}
		f, err := schemaDir.Open(path)
		if err != nil {
			return err
		}
		compiler.AddResource(path, f)
		return err
	})
	if err != nil {
		log.Fatalf("error loading schemas: %v", err)
	}
}

type internalSchema struct {
	schemas map[string]string
	id      string
}
