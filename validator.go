package openapi

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/chanced/uri"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

//go:embed schema
var embeddedSchemas embed.FS

const (
	// URI for OpenAPI 3.1 schema
	OPEN_API_3_1_SCHEMA = "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
	// URI for OpenAPI 3.0 schema
	OPEN_API_3_0_SCHEMA = "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
	// URI for JSON Schema 2020-12
	JSON_SCHEMA_2020_12 = "https://json-schema.org/draft/2020-12/schema"
	// URI for JSON Schema 2019-09
	JSON_SCHEMA_2019_09 = "https://json-schema.org/draft/2019-09/schema"
)

var (
	// OpenAPI31Schema is the URI for the JSON Schema of OpenAPI 3.1
	OpenAPI31Schema = *uri.MustParse(OPEN_API_3_1_SCHEMA)
	// OpenAPI30Schema is the URI for the JSON Schema of OpenAPI 3.0
	OpenAPI30Schema = *uri.MustParse(OPEN_API_3_0_SCHEMA)
	// JSONSchema202012SchemaURI is the URI for JSON Schema 2020-12
	JSONSchemaDialect202012 = *uri.MustParse(JSON_SCHEMA_2020_12)
	// JSONSchemaDialect201909 is the URI for JSON Schema 2019-09
	JSONSchemaDialect201909 = *uri.MustParse(JSON_SCHEMA_2019_09)
	// VersionConstraints3_0 is a semantic versioning constraint for 3.0:
	//	>= 3.0.0, < 3.1.0
	VersionConstraints3_0 = mustParseConstraints(">= 3.0.0, < 3.1.0")
	// SemanticVersion3_0 is a semantic versioning constraint for 3.1:
	//	>= 3.1.0, < 3.2.0
	VersionConstraints3_1 = mustParseConstraints(">= 3.1.0, < 3.2.0")
	// SupportedVersions is a semantic versioning constraint for versions
	// supported by openapi
	//
	// This is currently:
	//	>= 3.0.0, < 3.2.0
	SupportedVersions = mustParseConstraints(">= 3.0.0, < 3.2.0")
	// Version3_1 is a semantic version for 3.1.x
	Version3_1 = *semver.MustParse("3.1")
	// Version3_0 is a semantic version for 3.0.x
	Version3_0 = *semver.MustParse("3.0")

	// // JSONSchemaDialect07 is the URI for JSON Schema 07
	// JSONSchemaDialect07 = *uri.MustParse("http://json-schema.org/draft-07/schema#")
	// // JSONSchemaDialect04 is the URI for JSON Schema 04
	// JSONSchemaDialect04 = *uri.MustParse("http://json-schema.org/draft-04/schema#")
)

var _ Validator = (*StdValidator)(nil)

type Validator interface {
	// Validate should validate the fully-resolved OpenAPI document.
	ValidateDocument(document *Document) error

	// ValidateComponent should validate the structural integrity of a of an OpenAPI
	// document or component.
	//
	// If $ref is present in the data and the data is not a Schema, the Kind will be KindReference.
	// Otherwise, it will be the Kind of the data being loaded.
	//
	// openapi should only ever call Validate for the following:
	//    - OpenAPI Document (KindDocument)
	//    - JSON Schema (KindSchema)
	//    - Components (KindComponents)
	//    - Callbacks (KindCallbacks)
	//    - Example (KindExample)
	//    - Header (KindHeader)
	//    - Link (KindLink)
	//    - Parameter (KindParameter)
	//    - PathItem (KindPathItem)
	//    - Operation (KindOperation)
	//    - Reference (KindReference)
	//    - RequestBody (KindRequestBody)
	//    - Response (KindResponse)
	//    - SecurityScheme (KindSecurityScheme)
	//
	// StdComponentValidator will return an error if CompiledSchemas does not contain
	// a CompiledSchema for the given Kind.
	Validate(data []byte, resource uri.URI, kind Kind, openapi semver.Version, jsonschema uri.URI) error
}

// NewStdValidator creates and returns a new StdValidator.
//
// compiler is used to compile JSON Schema for initial validation.

// Each fs.FS in resources will be walked and all files ending in .json will be
// be added to the compiler. Defaults are provided from an embedded fs.FS.
//
// ## Resource Defaults
//   - OpenAPI 3.1: "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
//   - OpenAPI 3.0: "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
//   - JSON Schema 2020-12: "https://json-schema.org/draft/2020-12/schema"
//   - JSON Schema 2019-09: "https://json-schema.org/draft/2019-09/schema"
func NewValidator(compiler *jsonschema.Compiler, resources ...fs.FS) (*StdValidator, error) {
	if compiler == nil {
		return nil, errors.New("openapi: compiler is required")
	}
	compiled, err := CompileSchemas(compiler)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schemas: %w", err)
	}
	return &StdValidator{
		Schemas: compiled,
	}, nil
}

// StdValidator is an implemtation of the Validator interface.
type StdValidator struct {
	Schemas CompiledSchemas
}

// Validate should validate the fully-resolved OpenAPI document.
//
// This currently only validates with JSON Schema.
func (sv *StdValidator) ValidateDocument(doc *Document) error {
	// The openapi spec claims there are validations which json
	// schema can not fully encompass. Those will need to be added here.
	// TODO: Improve validation beyond JSON Schema

	d, err := doc.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	dialect := doc.JSONSchemaDialect
	if dialect == nil {
		if VersionConstraints3_1.Check(doc.OpenAPI) {
			dialect = &JSONSchemaDialect202012
			// } else if VersionConstraints3_0.Check(doc.OpenAPI) {
			// 	dialect = &JSONSchemaDialect201909
		} else {
			return fmt.Errorf("openapi: unable to detect OpenAPI version: %s", doc.OpenAPI)
		}
	}
	if err = sv.Validate(d, doc.AbsoluteLocation(), KindDocument, *doc.OpenAPI, *dialect); err != nil {
		return err
	}
	m := map[string]struct{}{}

	for _, r := range doc.Refs() {
		u := r.URI()
		if _, ok := m[r.ResolvedNode().AbsoluteLocation().String()]; ok {
			continue
		} else {
			m[r.ResolvedNode().AbsoluteLocation().String()] = struct{}{}
		}
		if u.Path != "" || u.Host != "" {
			if s, ok := r.ResolvedNode().(*Schema); ok {
				sd := dialect
				if s.Schema != nil {
					sd = s.Schema
				}
				d, err := s.MarshalJSON()
				if err != nil {
					return fmt.Errorf("failed to marshal schema: %w", err)
				}
				if err = sv.Validate(d, s.AbsoluteLocation(), KindSchema, *doc.OpenAPI, *sd); err != nil {
					return err
				}
			} else {
				d, err := r.ResolvedNode().MarshalJSON()
				if err != nil {
					return fmt.Errorf("failed to marshal node: %w", err)
				}

				if err = sv.Validate(d, r.ResolvedNode().AbsoluteLocation(), r.ResolvedNode().Kind(), *doc.OpenAPI, *dialect); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sv *StdValidator) Validate(data []byte, resource uri.URI, kind Kind, openapi semver.Version, jsonschema uri.URI) error {
	var i interface{}

	if kind == KindSchema {
		schema, ok := sv.Schemas.JSONSchema[jsonschema]
		if !ok {
			return fmt.Errorf("openapi: no schema found for %q", jsonschema)
		}
		if err := json.Unmarshal(data, &i); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		return schema.Validate(i)
	}
	var s CompiledSchema
	var ok bool
	if VersionConstraints3_1.Check(&openapi) {
		s, ok = sv.Schemas.OpenAPI[Version3_1][kind]
	}

	if !ok {
		return fmt.Errorf("openapi: schema not found for %s", kind)
	}

	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	if err := s.Validate(i); err != nil {
		return NewValidationError(err, kind, resource)
	}
	return nil
}

// CompiledSchema is an interface satisfied by a JSON Schema implementation	that
// validates primitive interface{} types.
//
// github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
type CompiledSchema interface {
	Validate(data interface{}) error
}

// // Compiler is an interface satisfied by any type which manages and compiles
// // resources (received in the form of io.Reader) based off of a URIs (including
// // fragments).
// //
// // github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
// type Compiler interface {
// 	AddResource(id string, r io.Reader) error
// 	Compile(url string) (CompiledSchema, error)
// }

// CompiledSchemas are used in the the StdValidator
type CompiledSchemas struct {
	OpenAPI    map[semver.Version]map[Kind]CompiledSchema
	JSONSchema map[uri.URI]CompiledSchema
}

// SetupCompiler adds OpenAPI and JSON Schema resources to a Compiler.
//
// Each fs.FS in resources will be walked and all files ending in .json will be
// be added to the compiler.
//
// # Defaults
//   - OpenAPI 3.1: "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
//   - OpenAPI 3.0: "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
//   - JSON Schema 2020-12: "https://json-schema.org/draft/2020-12/schema"
//   - JSON Schema 2019-09: "https://json-schema.org/draft/2019-09/schema"
func SetupCompiler(compiler *jsonschema.Compiler, resources ...fs.FS) (*jsonschema.Compiler, error) {
	if compiler == nil {
		return nil, errors.New("openapi: compiler is required")
	}
	resources = append([]fs.FS{embeddedSchemas}, resources...)
	err := addCompilerResources(compiler, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to add resources to compiler: %w", err)
	}
	return compiler, nil
}

// addCompilerResources adds the following schemas to a compiler:
//   - OpenAPI 3.1 ("https://spec.openapis.org/oas/3.1/schema/2022-02-27)")
//   - OpenAPI 3.0 ("https://spec.openapis.org/oas/3.0/schema/2021-09-28")
//   - JSON Schema 2020-12
//   - JSON Schema 2019-09
func addCompilerResources(compiler *jsonschema.Compiler, dirs []fs.FS) error {
	var err error
	for _, dir := range dirs {
		err = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".json" {
				return nil
			}
			f, err := dir.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()

			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			id := gjson.GetBytes(b, "$id").String()
			if id == "" {
				id = gjson.GetBytes(b, "id").String()
			}
			if id == "" {
				return fmt.Errorf("openapi: $id, id not found in %s", path)
			}
			err = compiler.AddResource(id, bytes.NewReader(b))
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func openAPISchemaVMap(openAPISchemas []map[string]uri.URI) (map[semver.Version]uri.URI, error) {
	res := make(map[semver.Version]uri.URI, 2)
	for _, vers := range openAPISchemas {
		for k, v := range vers {
			vers, err := semver.NewVersion(k)
			if err != nil {
				return nil, fmt.Errorf("failed to parse openAPISchemaID version: %w", err)
			}

			if _, errs := SupportedVersions.Validate(vers); err != nil {
				return nil, &UnsupportedVersionError{Version: k, Errs: errs}
			}

			k = fmt.Sprintf("%d.%d", vers.Major(), vers.Minor())
			vers, _ = semver.NewVersion(k)
			res[*vers] = v
		}
	}
	if _, ok := res[Version3_1]; !ok {
		res[Version3_1] = OpenAPI31Schema
	}
	if _, ok := res[Version3_0]; !ok {
		res[Version3_0] = OpenAPI30Schema
	}
	return res, nil
}

// CompileSchemas compiles the OpenAPI and JSON Schema resources using compiler.
//
// openAPISchemas is a variadic map of OpenAPI versions to their
// respective schema ids. The keys must be valid semver versions; only the major
// and minor versions are used. The last value for a given major and minor will be
// used
//
// Default openAPISchemas:
//
//	{ "3.1": "https://spec.openapis.org/oas/3.1/schema/2022-02-27)" }
//	{ "3.0": "https://spec.openapis.org/oas/3.0/schema/2021-09-28"  }
func CompileSchemas(compiler *jsonschema.Compiler, openAPISchemas ...map[string]uri.URI) (CompiledSchemas, error) {
	var err error

	openapis, err := compileOpenAPISchemas(compiler, openAPISchemas)
	if err != nil {
		return CompiledSchemas{}, fmt.Errorf("openapi: failed to compile openAPISchemas: %w", err)
	}
	jsonschemas, err := compileJSONSchemaSchemas(compiler)
	if err != nil {
		return CompiledSchemas{}, fmt.Errorf("openapi: failed to compile jsonschema schemas: %w", err)
	}
	return CompiledSchemas{
		OpenAPI:    openapis,
		JSONSchema: jsonschemas,
	}, nil
}

func compileJSONSchemaSchemas(c *jsonschema.Compiler) (map[uri.URI]CompiledSchema, error) {
	var err error
	jsonschemas := make(map[uri.URI]CompiledSchema, 2)
	jsonschemas[JSONSchemaDialect202012], err = c.Compile(JSON_SCHEMA_2020_12)
	if err != nil {
		return nil, err
	}
	jsonschemas[JSONSchemaDialect201909], err = c.Compile(JSON_SCHEMA_2019_09)
	if err != nil {
		return nil, err
	}
	return jsonschemas, nil
}

func compileOpenAPISchemas(c *jsonschema.Compiler, openAPISchemas []map[string]uri.URI) (map[semver.Version]map[Kind]CompiledSchema, error) {
	vm, err := openAPISchemaVMap(openAPISchemas)
	if err != nil {
		return nil, err
	}
	compiled := make(map[semver.Version]map[Kind]CompiledSchema, len(vm))

	for k, v := range vm {
		compiled[k], err = compileOpenAPISchemasFor(c, v)
		if err != nil {
			return nil, fmt.Errorf("openapi: failed to compile OpenAPI %s Schema %s: %w", k.String(), v.String(), err)
		}
	}
	return compiled, nil
}

func compileOpenAPISchemasFor(compiler *jsonschema.Compiler, uri uri.URI) (map[Kind]CompiledSchema, error) {
	uri.Fragment = ""
	uri.RawFragment = ""
	spec := uri.String()
	compileDef := func(name string) (CompiledSchema, error) {
		if uri.String() == "https://spec.openapis.org/oas/3.0/schema/2021-09-28" {
			if name == "callbacks" {
				name = "paths"
			}
			return compiler.Compile(spec + "#/definitions/" + Text(name).ToCamel().String())
		} else {
			return compiler.Compile(spec + "#/$defs/" + name)
		}
	}

	document, err := compiler.Compile(spec)
	if err != nil {
		return nil, fmt.Errorf("error compiling Dcument schema: %w", err)
	}

	operation, err := compileDef("operation")
	if err != nil {
		return nil, fmt.Errorf("error compiling Operation schema: %w", err)
	}

	callbacks, err := compileDef("callbacks")
	if err != nil {
		return nil, fmt.Errorf("error compiling Callbacks schema: %w", err)
	}

	example, err := compileDef("example")
	if err != nil {
		return nil, fmt.Errorf("error compiling Example schema: %w", err)
	}

	header, err := compileDef("header")
	if err != nil {
		return nil, fmt.Errorf("error compiling Header schema: %w", err)
	}

	link, err := compileDef("link")
	if err != nil {
		return nil, fmt.Errorf("error compiling Link schema: %w", err)
	}

	parameter, err := compileDef("parameter")
	if err != nil {
		return nil, fmt.Errorf("error compiling Parameter schema: %w", err)
	}

	requestBody, err := compileDef("request-body")
	if err != nil {
		return nil, fmt.Errorf("error compiling RequestBody schema: %w", err)
	}

	pathItem, err := compileDef("path-item")
	if err != nil {
		return nil, fmt.Errorf("error compiling PathItem schema: %w", err)
	}

	response, err := compileDef("response")
	if err != nil {
		return nil, fmt.Errorf("error compiling Response schema: %w", err)
	}

	securityScheme, err := compileDef("security-scheme")
	if err != nil {
		return nil, fmt.Errorf("error compiling SecurityScheme schema: %w", err)
	}
	reference, err := compileDef("reference")
	if err != nil {
		return nil, fmt.Errorf("error compiling Reference schema: %w", err)
	}

	o := map[Kind]CompiledSchema{
		KindDocument:       document,
		KindOperation:      operation,
		KindCallbacks:      callbacks,
		KindExample:        example,
		KindHeader:         header,
		KindLink:           link,
		KindParameter:      parameter,
		KindRequestBody:    requestBody,
		KindResponse:       response,
		KindSecurityScheme: securityScheme,
		KindPathItem:       pathItem,
		KindReference:      reference,
	}
	return o, nil
}

func mustParseConstraints(str string) semver.Constraints {
	c, err := semver.NewConstraint(str)
	if err != nil {
		log.Fatalf("failed to parse semver constraint %s: %v", str, err)
	}
	return *c
}
