package openapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"unicode"

	"github.com/chanced/openapi/yamlutil"
	jsonptr "github.com/xeipuuv/gojsonpointer"
)

// Opener is implemented by any value that has an Open method, which accepts a
// path (string) and returns an io.ReadCloser.
//
// The path provided will be in the form of "./path-to-resource", relative to
// the key provided to the Openers map.
//
// Two Openers are provided:
// 	- openapi.FSOpener: opens from an fs.FS
//	- openapi.HTTPOpener: opens by making HTTP requests
//
// Example:
//
//	import (
//		"log"
//		"embed"
//		"github.com/chanced/openapi"
//	)
//	//go:embed "openapi"
//	var embeddedfs embed.FS
//
//	oai, err := embeddedfs.Open("openapi.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
// 	o, err := openapi.Load(oai, openapi.NewResolver(openapi.Openers{
// 		"https://network.local": &openapi.FSOpener{FS: embeddedfs},
// 		"https://example.com": &openapi.HTTPOpener{},
//	}))
type Opener interface {
	Open(path string) (io.ReadCloser, error)
}

type openiniter interface {
	Init(string) error
}

// Openers is a map of URI addresses to Openers.
type Openers map[string]Opener

// ParameterResolverFunc resolves Parameters
type ParameterResolverFunc func(ref string) (*ParameterObj, error)

// ResponseResolverFunc resolves Responses
type ResponseResolverFunc func(ref string) (*ResponseObj, error)

// ExampleResolverFunc resolves Examples
type ExampleResolverFunc func(ref string) (*ExampleObj, error)

// HeaderResolverFunc resolves Headers
type HeaderResolverFunc func(ref string) (*HeaderObj, error)

// RequestBodyResolverFunc resolves RequestBodies
type RequestBodyResolverFunc func(ref string) (*RequestBodyObj, error)

// CallbackResolverFunc resolves Callbacks
type CallbackResolverFunc func(ref string) (*CallbackObj, error)

// PathResolverFunc resolves Paths
type PathResolverFunc func(ref string) (*PathObj, error)

// SecuritySchemeResolverFunc resolves SecuritySchemes
type SecuritySchemeResolverFunc func(ref string) (*SecuritySchemeObj, error)

// LinkResolverFunc resolves Links
type LinkResolverFunc func(ref string) (*LinkObj, error)

// SchemaResolverFunc resolves Schemas
type SchemaResolverFunc func(ref string) (*SchemaObj, error)

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

type resolver struct {
	openers map[string]Opener
	cache   *cache
}

func NewResolver(openers Openers) *resolver {
	for k, o := range openers {
		if oi, ok := o.(openiniter); ok {
			// ignoring errors; presumably they'll be returned by o.Open
			_ = oi.Init(k)
		}
	}
	dr := &resolver{
		openers: openers,
		cache:   newCache(),
	}
	return dr
}

type readercloser struct {
	io.Reader
	io.Closer
}

func (cr *resolver) opener(p string) (string, Opener, error) {
	if p == "" {
		return "", nil, errors.New("openapi: ref must not be empty")
	}
	for k, o := range cr.openers {
		if strings.HasPrefix(p, k) {
			return k, o, nil
		}
	}
	return "", nil, errors.New("openapi: no opener for " + p)
}

func (cr *resolver) open(pth string) (io.ReadCloser, error) {
	u, o, err := cr.opener(pth)
	if err != nil {
		return nil, err
	}
	pth = removeURI(pth, u)
	var ptr string
	pth, ptr = splitRef(pth)
	rc, err := o.Open(pth)
	if err != nil {
		return nil, err
	}
	if ptr != "" {
		return readPtr(rc, ptr)
	}
	return rc, nil
}
func (cr *resolver) ResolveParameterResolver(ref string) (*ParameterObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveResponseResolver(p string) (*ResponseObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveExampleResolver(p string) (*ExampleObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveHeaderResolver(p string) (*HeaderObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveRequestBodyResolver(p string) (*RequestBodyObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveCallbackResolver(p string) (*CallbackObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolvePathResolver(p string) (*PathObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveSecuritySchemeResolver(p string) (*SecuritySchemeObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveLinkResolver(p string) (*LinkObj, error) {
	panic("not implemented") // TODO: Implement
}

func (cr *resolver) ResolveSchemaResolver(p string) (*SchemaObj, error) {
	panic("not implemented") // TODO: Implement
}

type cache struct {
	Params          map[string]*ParameterObj
	Responses       map[string]*ResponseObj
	Examples        map[string]*ExampleObj
	Headers         map[string]*HeaderObj
	RequestBodies   map[string]*RequestBodyObj
	Callbacks       map[string]*CallbackObj
	Paths           map[string]*PathObj
	SecuritySchemes map[string]*SecuritySchemeObj
	Links           map[string]*LinkObj
	Schemas         map[string]*SchemaObj
}

func newCache() *cache {
	return &cache{
		Params:          make(map[string]*ParameterObj),
		Responses:       make(map[string]*ResponseObj),
		Examples:        make(map[string]*ExampleObj),
		Headers:         make(map[string]*HeaderObj),
		RequestBodies:   make(map[string]*RequestBodyObj),
		Callbacks:       make(map[string]*CallbackObj),
		Paths:           make(map[string]*PathObj),
		SecuritySchemes: make(map[string]*SecuritySchemeObj),
		Links:           make(map[string]*LinkObj),
		Schemas:         make(map[string]*SchemaObj),
	}
}

// FSOpener opens files from a filesystem
type FSOpener struct {
	FS fs.FS
}

func (o FSOpener) Open(name string) (io.ReadCloser, error) {
	return o.FS.Open(name)
}

// HTTPOpener opens files by making HTTPOpener requests
type HTTPOpener struct {
	URL            string
	url            *url.URL
	Client         *http.Client
	PrepareRequest func(req http.Request) http.Request
}

func (o *HTTPOpener) Init(v string) error {
	if v == "" {
		return errors.New("openapi: HTTPOpener URL must not be empty")
	}
	if o.URL == "" {
		o.URL = v
	}
	if o.url == nil {
		u, err := url.Parse(o.URL)
		if err == nil {
			o.url = u
		}
		return err
	}
	if o.Client == nil {
		o.Client = &http.Client{}
	}
	return nil
}

// Open opens the remote JSON or YAML by making an HTTP GET request, returning a io.Reader
func (o *HTTPOpener) Open(name string) (io.ReadCloser, error) {
	if o.url == nil || o.Client == nil {
		if err := o.Init(o.URL); err != nil {
			return nil, err
		}
	}
	if o.Client == nil {
		o.Client = &http.Client{}
	}
	addr := path.Join(o.url.Path, name)
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}
	if o.PrepareRequest != nil {
		mr := o.PrepareRequest(*req)
		req = &mr
	}
	res, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, err
}

const (
	jsonEncoding = iota + 1
	yamlEncoding
)

func decodePtr(dec *json.Decoder, ptr string, dst interface{}) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(dst)}
	}
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return err
	}
	p, err := jsonptr.NewJsonPointer(ptr)
	if err != nil {
		return err
	}
	pv, _, err := p.Get(&v)
	if err != nil {
		return err
	}
	b, err := json.Marshal(pv)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func decode(r io.Reader, ptr string, dst interface{}) error {
	var err error
	var enc uint8
	r, enc, err = detectEncoding(r)
	if err != nil {
		return err
	}
	if enc == yamlEncoding {
		r, err = yamlutil.EncodeYAMLToJSON(r)
		if err != nil {
			return err
		}
	}
	d := json.NewDecoder(r)
	if len(ptr) > 0 {
		return decodePtr(d, ptr, dst)
	}
	return d.Decode(dst)
}

func decodeAndClose(r io.ReadCloser, ptr string, dst interface{}) error {
	defer r.Close()
	return decode(r, ptr, dst)
}

func detectEncoding(r io.Reader) (io.Reader, uint8, error) {
	b := make([]byte, 1)
	var err error
	for {
		if _, err = r.Read(b); err != nil {
			if errors.Is(err, io.EOF) {
				return nil, 0, errors.New("unexpected EOF")
			}
			return nil, 0, err
		}
		if b[0] == '{' || b[0] == '[' {
			return io.MultiReader(bytes.NewReader(b), r), jsonEncoding, nil
		}
		if !unicode.IsSpace(rune(b[0])) {
			return io.MultiReader(bytes.NewReader(b), r), yamlEncoding, nil
		}
	}
}
func removeURI(r string, u string) string {
	r = strings.TrimPrefix(r, u)
	if r[0] == '/' {
		return "." + r
	}
	if r[0] == '#' {
		return r
	}

	return "./" + r
}

func splitRef(ref string) (string, string) {
	if len(ref) == 0 {
		return "", ""
	}
	if ref[0] == '#' {
		return "", ref[1:]
	}
	s := strings.SplitN(ref, "#", 2)
	if len(s) == 2 {
		if len(s[0]) > 0 && s[0][0] != '/' {
			s[0] = "/" + s[0]
		}
		return s[0], s[1]
	}
	return s[0], ""
}

func readPtr(rc io.ReadCloser, ptr string) (io.ReadCloser, error) {
	p, err := jsonptr.NewJsonPointer(ptr)
	if err != nil {
		return nil, err
	}
	r, e, err := detectEncoding(rc)
	if err != nil {
		return nil, err
	}
	var b []byte
	if e == yamlEncoding {
		b, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b, err = yamlutil.YAMLToJSON(b)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(b)
	}
	var c interface{}
	if err = json.NewDecoder(r).Decode(&c); err != nil {
		return nil, err
	}
	v, _, err := p.Get(c)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err = enc.Encode(v); err != nil {
		return nil, err
	}
	return &readercloser{
		Reader: buf,
		Closer: rc,
	}, nil
}
