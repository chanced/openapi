package openapi

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/chanced/uri"
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
	OpenAPI31SchemaURI          = *uri.MustParse(OPEN_API_3_1_SCHEMA)
	OpenAPI30SchemaURI          = *uri.MustParse(OPEN_API_3_0_SCHEMA)
	JSONSchema202012URI         = *uri.MustParse(JSON_SCHEMA_2020_12)
	JSONSchema201909URI         = *uri.MustParse(JSON_SCHEMA_2019_09)
	versionThreeOneConstraints  = mustParseConstraints(">=3.1.0 < 3.2.0")
	versionThreeZeroConstraints = mustParseConstraints(">=3.0.0 < 3.1.0")
	supportedVersions           = mustParseConstraints(">=3.0.0 < 3.2.0")
	v31                         = *semver.MustParse("3.1")
	v30                         = *semver.MustParse("3.0")
)

var _ Validator = (*StdValidator)(nil)

type Validator interface {
	// Validate should validate the OpenAPI document. It will be fully-resolved.
	//
	// StdValidator does not currently perform additional validation.
	ValidateDocument(document *Document) error

	// ValidateSchema should validate a JSON Schema document.
	//
	// dialect will be (in this order):
	//  - value of $schema, if present
	//  - value of jsonSchemaDialect if present in the nearest Document
	//  - "https://json-schema.org/draft/2020-12/schema" if OpenAPI is v3.1
	//  - "https://json-schema.org/draft/2019-09/schema" if OpenAPI is v3.0
	ValidateSchema(data []byte, uri *uri.URI, dialect *uri.URI) error

	// Validate should validate the structural integrity of a of an OpenAPI
	// document or component.
	//
	// If $ref is present in the data, the Kind will be KindReference.
	// Otherwise, it will be the Kind of the data being loaded.
	//
	// This should only be called for the following:
	//
	//    - Document (KindDocument)
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
	// StdValidator will return an error if the CompiledSchemas does not contani
	// a CompiledSchema for the given Kind.
	//
	// uri is the URI of the data being validated. It is used for error reporting
	Validate(data []byte, kind Kind, openapi semver.Version, uri *uri.URI) error
}

// NewStdValidator creates and returns a new StdValidator.
//
// compiler is used to compile JSON Schema for initial validation.
// The interface is satisfied by github.com/santhosh-tekuri/jsonschema/v5
//
// Each fs.FS in resources will be walked and all files ending in .json will be
// be added to the compiler. Defaults are provided from an embedded fs.FS.
//
// ## Resource Defaults
//   - OpenAPI 3.1: "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
//   - OpenAPI 3.0: "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
//   - JSON Schema 2020-12: "https://json-schema.org/draft/2020-12/schema"
//   - JSON Schema 2019-09: "https://json-schema.org/draft/2019-09/schema"
func NewValidator(compiler Compiler, resources ...fs.FS) (*StdValidator, error) {
	if compiler == nil {
		return nil, errors.New("openapi: compiler is required")
	}
	err := AddCompilerResources(compiler, resources...)
	if err != nil {
		return nil, fmt.Errorf("failed to add resources to compiler: %w", err)
	}
	compiled, err := CompileSchemas(compiler)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schemas: %w", err)
	}
	return &StdValidator{
		schemas: compiled,
	}, nil
}

// StdValidator is an implemtation of the Validator interface.
type StdValidator struct {
	schemas CompiledSchemas
}

// Validate implements Validator
func (sv *StdValidator) Validate(data []byte, kind Kind, _ *uri.URI, openapi semver.Version) error {
	if isRefJSON(data) {
		kind = KindReference
	}
	switch kind {
	case KindSchema:

		return sv.ValidateSchema(data)
	}
}

// ValidateDocument implements Validator
func (*StdValidator) ValidateDocument(document *Document) error {
	panic("unimplemented")
}

// ValidateSchema implements Validator
func (*StdValidator) ValidateSchema(data []byte, uri *uri.URI, dialect *uri.URI) error {
	panic("unimplemented")
}

// CompiledSchema is an interface satisfied by a JSON Schema implementation	that
// validates primitive interface{} types.
//
// github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
type CompiledSchema interface {
	Validate(data interface{}) error
}

// Compiler is an interface satisfied by any type which manages and compiles
// resources (received in the form of io.Reader) based off of a URIs (including
// fragments).
//
// github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
type Compiler interface {
	AddResource(id string, r io.Reader) error
	Compile(url string) (CompiledSchema, error)
}

// JSONSchemaResources used in the StdValidator
type JSONSchemaResources struct {
	Version string // uri of the json schema
	Schema  CompiledSchema
}

// CompiledSchemas are used in the the StdValidator
type CompiledSchemas struct {
	OpenAPI    map[semver.Version]map[Kind]CompiledSchema
	JSONSchema map[uri.URI]CompiledSchema
}

// AddCompilerResources adds OpenAPI and JSON Schema resources to a Compiler.
//
// Each fs.FS in resources will be walked and all files ending in .json will be
// be added to the compiler.
//
// # Defaults
//   - OpenAPI 3.1: "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
//   - OpenAPI 3.0: "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
//   - JSON Schema 2020-12: "https://json-schema.org/draft/2020-12/schema"
//   - JSON Schema 2019-09: "https://json-schema.org/draft/2019-09/schema"
func AddCompilerResources(compiler Compiler, resources ...fs.FS) error {
	if compiler == nil {
		return errors.New("openapi: compiler is required")
	}
	resources = append([]fs.FS{embeddedSchemas}, resources...)
	return addCompilerResources(compiler, resources)
}

// addCompilerResources adds the following schemas to a compiler:
//   - OpenAPI 3.1 ("https://spec.openapis.org/oas/3.1/schema/2022-02-27)")
//   - OpenAPI 3.0 ("https://spec.openapis.org/oas/3.0/schema/2021-09-28")
//   - JSON Schema 2020-12
//   - JSON Schema 2019-09
func addCompilerResources(compiler Compiler, dirs []fs.FS) error {
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

			if _, errs := supportedVersions.Validate(vers); err != nil {
				return nil, &SupportedVersionError{Version: k, Errs: errs}
			}

			k = fmt.Sprintf("%d.%d", vers.Major(), vers.Minor())
			vers, err = semver.NewVersion(k)
			if err != nil {
				// this should never happen
				panic(fmt.Errorf("failed to parse semver %q: %v", k, err))
			}
			res[*vers] = v
		}
	}
	if _, ok := res[v31]; !ok {
		res[v31] = OpenAPI31SchemaURI
	}
	if _, ok := res[v30]; !ok {
		res[v30] = OpenAPI30SchemaURI
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
func CompileSchemas(compiler Compiler, openAPISchemas ...map[string]uri.URI) (CompiledSchemas, error) {
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

func compileJSONSchemaSchemas(c Compiler) (map[uri.URI]CompiledSchema, error) {
	var err error
	jsonschemas := make(map[uri.URI]CompiledSchema, 2)
	jsonschemas[JSONSchema202012URI], err = c.Compile(JSON_SCHEMA_2020_12)
	if err != nil {
		return nil, fmt.Errorf("openapi: failed to compile JSON Schema 2020-12: %w", err)
	}
	jsonschemas[JSONSchema201909URI], err = c.Compile(JSON_SCHEMA_2019_09)
	if err != nil {
		return nil, fmt.Errorf("openapi: failed to compile JSON Schema 2019-09: %w", err)
	}
	return jsonschemas, nil
}

func compileOpenAPISchemas(c Compiler, openAPISchemas []map[string]uri.URI) (map[semver.Version]map[Kind]CompiledSchema, error) {
	vm, err := openAPISchemaVMap(openAPISchemas)
	if err != nil {
		return nil, fmt.Errorf("openapi: failed to compile OpenAPI Schemas: %w", err)
	}
	compiled := make(map[semver.Version]map[Kind]CompiledSchema, len(vm))

	for k, v := range vm {
		compiled[k], err = compileOpenAPISchemasFor(c, v)
		if err != nil {
			return nil, fmt.Errorf("openapi: failed to compile OpenAPI %s Schema %s: %w", k, v, err)
		}
	}
	return compiled, nil
}

func compileOpenAPISchemasFor(compiler Compiler, uri uri.URI) (map[Kind]CompiledSchema, error) {
	// u := "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
	uri.Fragment = ""
	uri.RawFragment = ""
	spec := uri.String()
	compileDef := func(name string) (CompiledSchema, error) {
		return compiler.Compile(spec + "#/$defs/" + name)
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
