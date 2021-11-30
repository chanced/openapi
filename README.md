# openapi

Package openapi is a set of Go types for [OpenAPI Specification
3.1](https://spec.openapis.org/oas/v3.1.0). The primary purpose of the package
is to assist in generation of OpenAPI documentation or to offer building blocks
for code-generation.

## Documentation

[Documentation can be found on pkg.go.dev](https://pkg.go.dev/github.com/chanced/openapi).

## Validation

Currently, specifications are validated with JSON Schema. Per OpenAPI's
documentation, this may not be enough to properly encapsulate all the nuances
of a specification. However, JSON Schema is able to properly validate the current
OpenAPI 3.1 Specification test suite.

If you run into an edge case that is not validated adequetely, please open a ticket.

## Contributions

Please feel free to open up an issue or create a pull request if there are features
you'd like to contribute or issues.

## Dependencies

-   [github.com/santhosh-tekuri/jsonschema/v5](https://github.com/santhosh-tekuri/jsonschema/v5) (used for json schema validation)
-   [github.com/evanphx/json-patch/v5](https://github.com/evanphx/json-patch/v5) (used for testing purposes)
-   [github.com/stretchr/testify](https://github.com/stretchr/testify) (testing)
-   [github.com/tidwall/gjson](https://github.com/tidwall/gjson) (json parsing)
-   [github.com/tidwall/sjson](https://github.com/tidwall/sjson) (json manipulation)
-   [github.com/wI2L/jsondiff](https://github.com/wI2L/jsondiff) (testing purposes)
-   [gopkg.in/yaml.v2](https://github.com/wI2L/jsondiff) (yaml)
-   [sigs.k8s.io/yaml](https://sigs.k8s.io/yaml) (yaml)
-   [github.com/chanced/cmpjson](https://github.com/chanced/cmpjson) (testing purposes)
-   [github.com/chanced/dynamic](https://github.com/chanced/dynamic) (json parsing)
-   [github.com/pkg/errors](https://github.com/pkg/errors) (errors)

## License

MIT
