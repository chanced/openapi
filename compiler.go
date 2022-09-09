package openapi

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

//go:embed schema
var schemaDir embed.FS

type internalSchemas struct {
	schema202012 *jsonschema.Schema
	openapi3_1   openapi31Schemas
}

var schemas internalSchemas

func (s *internalSchemas) compile(compiler *jsonschema.Compiler) error {
	var err error
	s.schema202012, err = compiler.Compile("https://json-schema.org/draft/2020-12/schema")
	if err != nil {
		return err
	}
	return s.openapi3_1.compile(compiler)
}

type openapi31Schemas struct {
	OpenAPI        *jsonschema.Schema
	Operation      *jsonschema.Schema
	Callback       *jsonschema.Schema
	Example        *jsonschema.Schema
	Header         *jsonschema.Schema
	License        *jsonschema.Schema
	Link           *jsonschema.Schema
	Parameter      *jsonschema.Schema
	RequestBody    *jsonschema.Schema
	Response       *jsonschema.Schema
	SecurityScheme *jsonschema.Schema
	Tag            *jsonschema.Schema
	Paths          *jsonschema.Schema
	PathItem       *jsonschema.Schema
	MediaType      *jsonschema.Schema
	Info           *jsonschema.Schema
	Contact        *jsonschema.Schema
	Encoding       *jsonschema.Schema
	ExternalDocs   *jsonschema.Schema
	Reference      *jsonschema.Schema
	// XML          *jsonschema.Schema
}

func (oapi *openapi31Schemas) compile(compiler *jsonschema.Compiler) error {
	u := "https://spec.openapis.org/oas/3.1/schema/2022-02-27"

	compileDef := func(name string) (*jsonschema.Schema, error) {
		return compiler.Compile(u + "#/$defs/" + name)
	}
	openAPI, err := compiler.Compile(u)
	if err != nil {
		return err
	}

	operation, err := compileDef("operation")
	if err != nil {
		return err
	}

	callback, err := compileDef("callbacks")
	if err != nil {
		return err
	}

	example, err := compileDef("example")
	if err != nil {
		return err
	}

	header, err := compileDef("header")
	if err != nil {
		return err
	}

	license, err := compileDef("license")
	if err != nil {
		return err
	}

	link, err := compileDef("link")
	if err != nil {
		return err
	}

	parameter, err := compileDef("parameter")
	if err != nil {
		return err
	}

	requestBody, err := compileDef("request-body")
	if err != nil {
		return err
	}

	paths, err := compileDef("paths")
	if err != nil {
		return err
	}

	pathItem, err := compileDef("path-item")
	if err != nil {
		return err
	}

	response, err := compileDef("response")
	if err != nil {
		return err
	}

	securityScheme, err := compileDef("security-scheme")
	if err != nil {
		return err
	}

	tag, err := compileDef("tag")
	if err != nil {
		return err
	}

	mediaType, err := compileDef("media-type")
	if err != nil {
		return err
	}

	info, err := compileDef("info")
	if err != nil {
		return err
	}

	contact, err := compileDef("contact")
	if err != nil {
		return err
	}

	encoding, err := compileDef("encoding")
	if err != nil {
		return err
	}

	reference, err := compileDef("reference")
	if err != nil {
		return err
	}

	externalDocs, err := compileDef("external-documentation")
	if err != nil {
		return err
	}

	*oapi = openapi31Schemas{
		OpenAPI:        openAPI,
		Operation:      operation,
		Callback:       callback,
		Example:        example,
		Header:         header,
		License:        license,
		Link:           link,
		Parameter:      parameter,
		RequestBody:    requestBody,
		Response:       response,
		SecurityScheme: securityScheme,
		Tag:            tag,
		Paths:          paths,
		PathItem:       pathItem,
		MediaType:      mediaType,
		Info:           info,
		Contact:        contact,
		Encoding:       encoding,
		ExternalDocs:   externalDocs,
		Reference:      reference,
	}
	return nil
}

func init() {
	log.SetFlags(0)
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft2020
	err := fs.WalkDir(schemaDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		f, err := schemaDir.Open(path)
		if err != nil {
			return nil
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		id := gjson.GetBytes(b, "$id").String()
		err = compiler.AddResource(id, bytes.NewReader(b))
		return err
	})
	if err != nil {
		log.Fatalf("error loading schemas: %v", err)
	}
	err = schemas.compile(compiler)
	if err != nil {
		log.Fatalf("error compiling schemas: %v", err)
	}
}

type internalSchema struct {
	schemas map[string]string
	id      string
}
