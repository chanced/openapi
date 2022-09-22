# openapi - an OpenAPI 3.x library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/chanced/openapi.svg)](https://pkg.go.dev/github.com/chanced/openapi)

openapi is a library for parsing and validating OpenAPI
[3.1](https://spec.openapis.org/oas/v3.1.0),
[3.0](https://spec.openapis.org/oas/v3.0.3). The intent of the library is to
offer building blocks for code and documentation generation.

:warning: This library is under development; there may be breaking changes and bugs.

## Features

-   `$ref` resolution
-   All keys retain their order from the markup using slices of key/values which
    aids in code generation.
-   Validation ([see the validation section](#validation))
-   All non-primitive nodes have an absolute & relative location
-   Strings are [text.Text](https://github.com/chanced/caps) which has case
    conversions and `strings` functions as methods.
-   Extensions, unknown JSON Schema keywords, examples, and a few other fields
    are instances of [jsonx.RawMessage](https://github.com/chanced/jsonx) which
    comes with a few helper methods.
-   Supports both JSON and YAML

## Usage

```go
package main

import (
    "github.com/chanced/openapi"
    "github.com/chanced/uri"
    "github.com/santhosh-tekuri/jsonschema/v5"
    "embed"
    "io"
    "path/filepath"
    "log"
)

//go:embed spec
var specFiles embed.FS

func main() {
    ctx := context.Background()

    c, err := openapi.SetupCompiler(jsonschema.NewCompiler()) // adding schema files
    if err != nil {
        log.Fatal(err)
    }
    v, err := openapi.NewValidator(c)
    if err != nil {
        log.Fatal(err)
    }

    fn := func(_ context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error){
        f, err := specFiles.Open(fp)
        if err != nil {
            log.Fatal(err)
        }
        // you can return either JSON or YAML
        d, err := io.ReadAll(f)
        if err != nil{
            log.fatal(err)
        }
        // use the uri or the data to determine the Kind
        return openapi.KindDocument, d, nil
    }
    // you can Load either JSON or YAML
    // Load validates the Document as well.
    doc, err := openapi.Load(ctx, "spec/openapi.yaml", v, fn)
    if err != nil{
        log.Fatal(err)
    }
    _ = doc // *openapi.Document
}
```

## Validation

The standard validator (`StdValidator`) currently validates OpenAPI documents
with JSON Schema. Per OpenAPI's documentation, this may not be enough to
properly encapsulate all the nuances of a specification. However, JSON Schema is
able to successfully validate the current OpenAPI 3.1 Specification test suite.

Validation is an area that still needs work. If you do find cases where the
current validator is not sufficient, please open an issue so that the library
can be updated with proper coverage of that case.

## Contributions

Please feel free to open up an issue or create a pull request if you find issues
or if there are features you'd like to see added.

## License

MIT
