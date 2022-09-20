# openapi - an OpenAPI 3.x library for Go

<p align="center">
<b>:warning: This library is in an alpha state; expect breaking changes and bugs. :warning:</b>
</p>
openapi is a library for for OpenAPI 3.x, including
[3.1](https://spec.openapis.org/oas/v3.1.0) and
[3.0](https://spec.openapis.org/oas/v3.0.3). The primary purpose of the package
is to offer building blocks for code and documentation generation.

## Features

-   Reference resolution, including support for recursive `$ref`s, `$dynamicRef`,
    `$recursiveRef`.
-   Validation ([see the validation seciton](#validation))
-   All strings are instances of
    [github.com/chanced/caps/text](https://github.com/chanced/caps)
    which come with case conversion and `strings` functions as methods.
-   Extensions, unknown JSON Schema keywords, examples, and a few other things
    are instances of [jsonx.RawMessage](https://github.com/chanced/jsonx) which
    comes with a few helper methods.
-   All keys retain their order. This makes for a bit more of a complicated API
    but prevents the need for further resorting to keep code generation consistent

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
)

//go:embed spec
var spec embed.FS

func main() {
    ctx := context.Background()
    c, _ := openapi.SetupCompiler(jsonschema.NewCompiler())
    v := openapi.NewValidator(c)
    fn := func(_ context.Context, uri uri.URI, kind Kind) (Kind, []byte, error){
        // quick and simple example
        // be sure to handle errors
        f, _ := schema.Open(fp)
        // you can return either JSON or YAML
        d, _ := io.ReadAll(f)
        // use the uri or the data to determine the Kind
        return KindDocument, d, nil
    }
    // you can Load either JSON or YAML
    doc, _ := openapi.Load(ctx, "spec/openapi.yaml", v, fn)
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

## Known trouble spots

-   **Testing.** I simply do not have the time to properly test this at the moment.
    There be dragons :dragon_face: and bugs :lady_beetle:. If you find them, please [open an issue](https://github.com/chanced/openapi/issues/new).
-   **Validation.** [See the Validation section](#validation).
-   [jsonpointer](https://github.com/chanced/jsonpointer)'s Resolve, Assign, and
    Delete do not currently work. I need to update
    [jsonpointer](https://github.com/chanced/jsonpointer) before the interfaces
    can be implemented for types within this library.
-   Values of `$anchor` and `$dynamicAnchor` must be unique to a file.
    Conditional `$dynamicAnchor` `$recursiveAnchor` are going to be challenging.
    See below.
-   `$dynamicRef` and `$recursiveRef` are incredibly tricky with respect to
    static analysis, which is what this library was built for. You should avoid
    conditional branches with `$dynamicAnchor`s within the same file. If you
    need a conditional dynamics, move the branch into its own file and have the
    conditional statement reference the branch.

## Contributions

Please feel free to open up an issue or create a pull request if you find issues
or if there are features you'd like to see added.

## License

MIT
