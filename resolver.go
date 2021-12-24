package openapi

import (
	"errors"
	"io"
	"strings"
)

type JSONPointerResolver interface {
	ResolveJSONPointer(string, interface{}) error
}

// Resolver is implemented by any value that has ResolverFuns for each of the
// referencable OpenAPI objects
type Resolver interface {
	ResolveParameterResolver(string) (*ParameterObj, error)
	ResolveResponseResolver(string) (*ResponseObj, error)
	ResolveExampleResolver(string) (*ExampleObj, error)
	ResolveHeaderResolver(string) (*HeaderObj, error)
	ResolveRequestBodyResolver(string) (*RequestBodyObj, error)
	ResolveCallbackResolver(string) (*CallbackObj, error)
	ResolvePathResolver(string) (*PathObj, error)
	ResolveSecuritySchemeResolver(string) (*SecuritySchemeObj, error)
	ResolveLinkResolver(string) (*LinkObj, error)
	ResolveSchemaResolver(string) (*SchemaObj, error)
}

type OpenAPIResolver struct {
	openers map[string]Opener
	cache   *cache
	openapi *OpenAPI
}

// NewResolver returns a new OpenAPIResolver which implements Resolver.
func NewResolver(openers Openers) *OpenAPIResolver {
	for k, o := range openers {
		if oi, ok := o.(openiniter); ok {
			// ignoring errors; presumably they'll be returned by o.Open
			_ = oi.Init(k)
		}
	}
	dr := &OpenAPIResolver{
		openers: openers,
	}
	return dr
}

type readercloser struct {
	io.Reader
	io.Closer
}

func (oar *OpenAPIResolver) ResolveParameterResolver(ref string) (*ParameterObj, error) {
	if v, ok := oar.cache.Params[ref]; ok {
		return v, nil
	}

	if strings.HasPrefix(ref, "#/components/parameters/") {
		c, ok := oar.openapi.Components.Parameters[strings.TrimPrefix(ref, "#/components/parameters/")]
	}
}

func (oar *OpenAPIResolver) ResolveResponseResolver(p string) (*ResponseObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveExampleResolver(p string) (*ExampleObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveHeaderResolver(p string) (*HeaderObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveRequestBodyResolver(p string) (*RequestBodyObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveCallbackResolver(p string) (*CallbackObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolvePathResolver(p string) (*PathObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveSecuritySchemeResolver(p string) (*SecuritySchemeObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveLinkResolver(p string) (*LinkObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) ResolveSchemaResolver(p string) (*SchemaObj, error) {
	panic("not implemented") // TODO: Implement
}

func (oar *OpenAPIResolver) opener(p string) (string, Opener, error) {
	if p == "" {
		return "", nil, errors.New("openapi: ref must not be empty")
	}
	for k, o := range oar.openers {
		if strings.HasPrefix(p, k) {
			return k, o, nil
		}
	}
	return "", nil, errors.New("openapi: no opener for " + p)
}
