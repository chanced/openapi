package openapi

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

type Validator interface {
	Validate(data interface{}) error
}

type Compiler interface {
	AddResource(id string, r io.Reader) error
	Compile(url string) Validator
}

//go:embed schema
var schemaDir embed.FS

type internalSchemas struct {
	schema202012 *jsonschema.Schema
	openapi31    map[Kind]*jsonschema.Schema
}

func setupCompiler(compiler Compiler) (Compiler, error) {
	if compiler == nil {
		return nil, fmt.Errorf("error: compiler is nil")
	}
	err := addCompilerResources(compiler)
	if err != nil {
		return nil, err
	}
	return compiler, nil
}

// AddCompilerResources adds the schemas for OpenAPI 3.1 & JSON Schema 2020-12 to compiler
func addCompilerResources(compiler Compiler) error {
	return fs.WalkDir(schemaDir, ".", func(path string, d fs.DirEntry, err error) error {
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
}

func compileSchemas(compiler *jsonschema.Compiler) (internalSchemas, error) {
	var err error
	schemas := internalSchemas{}
	schemas.schema202012, err = compiler.Compile("https://json-schema.org/draft/2020-12/schema")
	if err != nil {
		return schemas, err
	}
	schemas.openapi31, err = compileOpenAPI31Schemas(compiler)
	return schemas, err
}

func compileOpenAPI31Schemas(compiler *jsonschema.Compiler) (map[Kind]*jsonschema.Schema, error) {
	u := "https://spec.openapis.org/oas/3.1/schema/2022-02-27"

	compileDef := func(name string) (*jsonschema.Schema, error) {
		return compiler.Compile(u + "#/$defs/" + name)
	}
	openAPI, err := compiler.Compile(u)
	if err != nil {
		return nil, err
	}

	operation, err := compileDef("operation")
	if err != nil {
		return nil, err
	}

	callbacks, err := compileDef("callbacks")
	if err != nil {
		return nil, err
	}

	example, err := compileDef("example")
	if err != nil {
		return nil, err
	}

	header, err := compileDef("header")
	if err != nil {
		return nil, err
	}

	license, err := compileDef("license")
	if err != nil {
		return nil, err
	}

	link, err := compileDef("link")
	if err != nil {
		return nil, err
	}

	parameter, err := compileDef("parameter")
	if err != nil {
		return nil, err
	}

	requestBody, err := compileDef("request-body")
	if err != nil {
		return nil, err
	}

	paths, err := compileDef("paths")
	if err != nil {
		return nil, err
	}

	pathItem, err := compileDef("path-item")
	if err != nil {
		return nil, err
	}

	response, err := compileDef("response")
	if err != nil {
		return nil, err
	}

	securityScheme, err := compileDef("security-scheme")
	if err != nil {
		return nil, err
	}

	tag, err := compileDef("tag")
	if err != nil {
		return nil, err
	}

	mediaType, err := compileDef("media-type")
	if err != nil {
		return nil, err
	}

	info, err := compileDef("info")
	if err != nil {
		return nil, err
	}

	contact, err := compileDef("contact")
	if err != nil {
		return nil, err
	}

	encoding, err := compileDef("encoding")
	if err != nil {
		return nil, err
	}

	reference, err := compileDef("reference")
	if err != nil {
		return nil, err
	}

	externalDocs, err := compileDef("external-documentation")
	if err != nil {
		return nil, err
	}

	o := map[Kind]*jsonschema.Schema{
		KindDocument:       openAPI,
		KindOperation:      operation,
		KindCallbacks:      callbacks,
		KindExample:        example,
		KindHeader:         header,
		KindLicense:        license,
		KindLink:           link,
		KindParameter:      parameter,
		KindRequestBody:    requestBody,
		KindResponse:       response,
		KindSecurityScheme: securityScheme,
		KindTag:            tag,
		KindPaths:          paths,
		KindPathItem:       pathItem,
		KindMediaType:      mediaType,
		KindInfo:           info,
		KindContact:        contact,
		KindEncoding:       encoding,
		KindExternalDocs:   externalDocs,
		KindReference:      reference,
	}
	return o, nil
}
