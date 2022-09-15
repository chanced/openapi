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
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

const (
	OPEN_API_3_1_SCHEMA = "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
	OPEN_API_3_0_SCHEMA = "https://spec.openapis.org/oas/3.0/schema/2021-09-28"
	JSON_SCHEMA_2020_12 = "https://json-schema.org/draft/2020-12/schema"
	JSON_SCHEMA_2019_09 = "https://json-schema.org/draft/2019-09/schema"
)

var (
	OpenAPI31SchemaURI          = *uri.MustParse(OPEN_API_3_1_SCHEMA)
	OpenAPI30SchemaURI          = *uri.MustParse(OPEN_API_3_0_SCHEMA)
	JSONSchema202012URI         = *uri.MustParse(JSON_SCHEMA_2020_12)
	JSONSchema201909URI         = *uri.MustParse(JSON_SCHEMA_2019_09)
	versionThreeOneConstraints  = mustParseConstraints(">=3.1.0 < 3.2.0")
	versionThreeZeroConstraints = mustParseConstraints(">=3.0.0 < 3.1.0")
	versionThreeOne             = *semver.MustParse("3.1.0")
	versionThreeZero            = *semver.MustParse("3.0.0")
)

//go:embed schema
var embeddedSchemas embed.FS

type Validator interface {
	// ValidateDocument should roughly validate data as an OpenAPI Document for
	// structural integerity. A second call to Validate will be made with the
	// fully resolved Document.
	ValidateDocument(data []byte, uri *uri.URI) error

	// Validate should validate the OpenAPI document. It will be fully-resolved.
	Validate(document *Document) error

	// ValidateSchema should validate a JSON Schema document.
	//
	// dialect will be (in this order):
	//  - value of $schema, if present
	//  - value of jsonSchemaDialect if present in the nearest Document
	//  - "https://json-schema.org/draft/2020-12/schema" if OpenAPI is v3.1
	//  - "https://json-schema.org/draft/2019-09/schema" if OpenAPI is v3.0
	ValidateSchema(data []byte, uri *uri.URI, dialect *uri.URI) error

	// ValidateComponent should validate the structural integrity of a
	// referenced component.
	//
	// If $ref is present in the data, the data will be validated against a
	// Reference. Otherwise the target component's Kind will be provided.
	//
	// This is only called for the following:
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
	ValidateComponent(data []byte, uri *uri.URI, openapi semver.Version) error

	// ValidateCallbacks(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateExample(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateHeader(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateLink(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateParameter(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidatePathItem(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateOperation(data []byte, uri *uri.URI, openapi semver.Version) error
	// ValidateReference(data []byte, uri *uri.URI, openapi*semver.Version) error
	// ValidateRequestBody(data []byte, uri *uri.URI, openapi *semver.Version) error
	// ValidateResponse(data []byte, uri *uri.URI, openapi *semver.Version) error
	// ValidateSecurityScheme(data []byte, uri *uri.URI, openapi *semver.Version) error
}

// CompiledSchema is an interface satisfied by a JSON Schema implementation	that
// validates primitive interface{} types.
//
// https://github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
type CompiledSchema interface {
	Validate(data interface{}) error
}

// Compiler is an interface satisfied by any type which manages and compiles
// resources (received in the form of io.Reader) based off of a URIs (including
// fragments).
//
// https://github.com/santhosh-tekuri/jsonschema/v5 satisfies this interface.
type Compiler interface {
	AddResource(id string, r io.Reader) error
	Compile(url string) (CompiledSchema, error)
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
		return errors.New("error: compiler is required")
	}
	resources = append([]fs.FS{embeddedSchemas}, resources...)
	return addCompilerResources(compiler, resources)
}

// CompiledOpenAPISchemas used in the StdValidator
type CompiledOpenAPISchemas struct {
	Version string // semver of the json schema
	Schemas map[Kind]CompiledSchema

	// Callbacks      CompiledSchema
	// Example        CompiledSchema
	// Header         CompiledSchema
	// Link           CompiledSchema
	// Parameter      CompiledSchema
	// PathItem       CompiledSchema
	// Operation      CompiledSchema
	// Reference      CompiledSchema
	// RequestBody    CompiledSchema
	// Response       CompiledSchema
	// SecurityScheme CompiledSchema
}

// JSONSchemaResources used in the StdValidator
type JSONSchemaResources struct {
	Version string // uri of the json schema
	Schema  CompiledSchema
}

// CompiledSchemas are used in the the StdValidator
type CompiledSchemas struct {
	OpenAPI map[semver.Version]CompiledOpenAPISchemas
	JSON    map[uri.URI]Compiler
}

type CompiledJSONSchema struct {
	Schema202012 Validator
	Schema201909 *jsonschema.Schema
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
	openapiURIs := make(map[string]uri.URI, 2)
	openapiVersLookup := make(map[string]semver.Version, 2)

	for _, vers := range openAPISchemas {
		for k, v := range vers {
			vers, err := semver.NewVersion(k)
			if err != nil {
				return CompiledSchemas{}, fmt.Errorf("error: failed to parse openAPISchemaID version: %w", err)
			}

			if vers.Major() != 3 {
				return CompiledSchemas{}, fmt.Errorf("error: OpenAPI version %s is not supported", vers)
			}
			k = fmt.Sprintf("%d.%d", vers.Major(), vers.Minor())
			openapiVersLookup[k] = *vers
			openapiURIs[k] = v
		}
	}

	if _, ok := openapiURIs["3.1"]; !ok {
		openapiURIs["3.1"] = OpenAPI31SchemaURI
		openapiVersLookup["3.1"] = versionThreeOne
	}
	if _, ok := openapiURIs["3.0"]; !ok {
		openapiURIs["3.0"] = OpenAPI30SchemaURI
		openapiVersLookup["3.0"] = versionThreeZero
	}

	var err error

	jsonschemas := make(map[uri.URI]CompiledSchema, 2)
	jsonschemas[JSONSchema202012URI], err = compiler.Compile("https://json-schema.org/draft/2020-12/schema")
	if err != nil {
		return CompiledSchemas{}, fmt.Errorf("error: failed to compile JSON Schema 2020-12: %w", err)
	}
	jsonschemas[JSONSchema201909URI], err = compiler.Compile("https://json-schema.org/draft/2019-09/schema")
	if err != nil {
		return CompiledSchemas{}, fmt.Errorf("error: failed to compile JSON Schema 2019-09: %w", err)
	}

	openapis := make(map[semver.Version]CompiledOpenAPISchemas, len(openapiURIs))
	_ = openapis
	// var openAPISchemas map[]
	var vers semver.Version
	for key, id := range openapiURIs {
		_ = id
		vers = openapiVersLookup[key]
		// openapis[vers], err = compiler.Compile(id.String())
		if err != nil {
			return CompiledSchemas{}, fmt.Errorf("error: failed to compile OpenAPI %s schema: %w", vers, err)
		}

	}
	panic("not done")
	// return schemas, err
}

func compileOpenAPI31Schemas(compiler Compiler, uri uri.URI) (map[Kind]CompiledSchema, error) {
	// u := "https://spec.openapis.org/oas/3.1/schema/2022-02-27"
	uri.Fragment = ""
	uri.RawFragment = ""
	spec := uri.String()
	compileDef := func(name string) (CompiledSchema, error) {
		return compiler.Compile(spec + "#/$defs/" + name)
	}

	document, err := compiler.Compile(spec)
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
	reference, err := compileDef("reference")
	if err != nil {
		return nil, err
	}

	// license, err := compileDef("license")
	// if err != nil {
	// 	return nil, err
	// }

	// paths, err := compileDef("paths")
	// if err != nil {
	// 	return nil, err
	// }

	// tag, err := compileDef("tag")
	// if err != nil {
	// 	return nil, err
	// }

	// mediaType, err := compileDef("media-type")
	// if err != nil {
	// 	return nil, err
	// }

	// info, err := compileDef("info")
	// if err != nil {
	// 	return nil, err
	// }

	// contact, err := compileDef("contact")
	// if err != nil {
	// 	return nil, err
	// }

	// encoding, err := compileDef("encoding")
	// if err != nil {
	// 	return nil, err
	// }

	// externalDocs, err := compileDef("external-documentation")
	// if err != nil {
	// 	return nil, err
	// }

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
		// KindLicense:        license,
		// KindTag:            tag,
		// KindPaths:          paths,
		// KindMediaType:      mediaType,
		// KindInfo:           info,
		// KindContact:        contact,
		// KindEncoding:       encoding,
		// KindExternalDocs:   externalDocs,

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
