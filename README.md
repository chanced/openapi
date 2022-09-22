# openapi - an OpenAPI 3.x library for Go

openapi is a library for for OpenAPI 3.x ([3.1](https://spec.openapis.org/oas/v3.1.0),
[3.0](https://spec.openapis.org/oas/v3.0.3)).

The primary purpose of the package is to offer building blocks for code and
documentation generation.

:warning: This library is in an alpha state; expect breaking changes and bugs.

## Features

-   `$ref` resolution
-   All keys retain their order from the markup using slices of key/values which
    aids with code generation.
-   Validation ([see the validation seciton](#validation))
-   All non-primitive nodes have an absolute & relative location
-   Strings are [text.Text](https://github.com/chanced/caps) which has case
    conversions and `strings` functions as methods.
-   Extensions, unknown JSON Schema keywords, examples, and a few other fields
    are instances of [jsonx.RawMessage](https://github.com/chanced/jsonx) which
    comes with a few helper methods.
-   Supports both JSON and YAML

## Issues

-   **Testing.** The code coverage is abysmal at the moment. As I find time, I'll add coverage.
-   **`$dynamicRef` / `$dynamicAnchor`** is not really supported. While the
    references are loaded, the dynamic overriding is not. I simply have no idea
    how to solve it. If you have ideas, I'd really like to hear them.
-   **Validation.** [See the Validation section](#validation).
-   **Errors.** Errors and error messages need a lot of work.
-   [jsonpointer](https://github.com/chanced/jsonpointer)'s Resolve, Assign, and
    Delete do not currently work. I need to update the jsonpointer library
    before its interfaces can be implemented for types within this library.
-   Values of `$anchor` and `$dynamicAnchor` must be unique to a file.
    Conditional `$dynamicAnchor` `$recursiveAnchor` are going to be challenging.
    See below.
-   `$dynamicRef` and `$recursiveRef` are incredibly tricky with respect to
    static analysis, which is what this library was built for. You should avoid
    conditional branches with `$dynamicAnchor`s within the same file. If you
    need a conditional dynamics, move the branch into its own file and have the
    conditional statement reference the branch.

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
var spec embed.FS

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
        f, err := schema.Open(fp)
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

Validation something that needs work. If you have an edge case that is not
covered, you can implement your own Validator either by wrapping `StdValidator`
or simply creating your own.

If you do find cases where the current validator is not sufficient, please open
an issue so that the library can be updated with proper coverage in the future.

Regarding JSON Schema, as of writing this, the only library able to support JSON
Schema 2020-12 is
[github.com/santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema)
and so the `Compiler`'s interface was modeled after its API. If you would like
to use a different implementation of JSON Schema with the `StdValidator` the
interfaces you need to write an adapter for are:

```go
type Compiler interface {
	AddResource(id string, r io.Reader) error
	Compile(url string) (CompiledSchema, error)
}

type CompiledSchema interface {
	Validate(data interface{}) error
}
```

## Contributions

Please feel free to open up an issue or create a pull request if you find issues
or if there are features you'd like to see added.

## License

MIT
