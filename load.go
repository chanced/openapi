package openapi

// // Load loads
// //
// // Example:
// //
// //	import (
// //		"log"
// //		"embed"
// //		"github.com/chanced/openapi"
// //	)
// //	//go:embed "openapi"
// //	var embeddedfs embed.FS
// //
// //	oai, err := embeddedfs.Open("openapi.yaml")
// //	if err != nil {
// //		log.Fatal(err)
// //	}
// //	o, err := openapi.Load(oai, openapi.NewResolver(openapi.Openers{
// //		"https://network.local": &openapi.FSOpener{FS: embeddedfs},
// //		"https://example.com": &openapi.HTTPOpener{},
// //	}))
// func Load(openapi io.Reader, resolver Resolver) (*ResolvedOpenAPI, error) {
// 	var o *OpenAPI
// 	if err := decode(openapi, "", &o); err != nil {
// 		return nil, err
// 	}
// 	return loader{cache: newCache(), resolver: resolver, openapi: o}.load()
// }

// type loader struct {
// 	*cache
// 	resolver Resolver
// 	openapi  *OpenAPI
// }

// func (l loader) load() (*ResolvedOpenAPI, error) {
// 	panic("not implemented") // TODO: implement
// }

// func (l loader) loadSchemas(schemas SchemaSet) {
// 	panic("not implemented") // TODO: implement
// }

// const (
// 	jsonEncoding = iota + 1
// 	yamlEncoding
// )

// func decodePtr(dec *json.Decoder, ptr string, dst interface{}) error {
// 	panic("not impl")
// 	// rv := reflect.ValueOf(dst)
// 	// if rv.Kind() != reflect.Ptr || rv.IsNil() {
// 	// 	return &json.InvalidUnmarshalError{Type: reflect.TypeOf(dst)}
// 	// }
// 	// var v interface{}
// 	// if err := dec.Decode(&v); err != nil {
// 	// 	return err
// 	// }
// 	// p, err := jsonptr.NewJsonPointer(ptr)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// pv, _, err := p.Get(&v)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// b, err := json.Marshal(pv)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// return json.Unmarshal(b, dst)
// }

// func decode(r io.Reader, ptr string, dst interface{}) error {
// 	var err error
// 	var enc uint8
// 	r, enc, err = detectEncoding(r)
// 	if err != nil {
// 		return err
// 	}
// 	if enc == yamlEncoding {
// 		r, err = yamlutil.EncodeYAMLToJSON(r)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	d := json.NewDecoder(r)
// 	if len(ptr) > 0 {
// 		return decodePtr(d, ptr, dst)
// 	}
// 	return d.Decode(dst)
// }

// func decodeAndClose(r io.ReadCloser, ptr string, dst interface{}) error {
// 	defer func(r io.ReadCloser) { _ = r.Close() }(r)
// 	return decode(r, ptr, dst)
// }

// func detectEncoding(r io.Reader) (io.Reader, uint8, error) {
// 	b := make([]byte, 1)
// 	var err error
// 	for {
// 		if _, err = r.Read(b); err != nil {
// 			if errors.Is(err, io.EOF) {
// 				return nil, 0, errors.New("unexpected EOF")
// 			}
// 			return nil, 0, err
// 		}
// 		if b[0] == '{' || b[0] == '[' {
// 			return io.MultiReader(bytes.NewReader(b), r), jsonEncoding, nil
// 		}
// 		if !unicode.IsSpace(rune(b[0])) {
// 			return io.MultiReader(bytes.NewReader(b), r), yamlEncoding, nil
// 		}
// 	}
// }

// func removeURI(r string, u string) string {
// 	r = strings.TrimPrefix(r, u)
// 	if r[0] == '/' {
// 		return "." + r
// 	}
// 	if r[0] == '#' {
// 		return r
// 	}

// 	return "./" + r
// }

// func splitRef(ref string) (string, string) {
// 	if len(ref) == 0 {
// 		return "", ""
// 	}
// 	if ref[0] == '#' {
// 		return "", ref[1:]
// 	}
// 	s := strings.SplitN(ref, "#", 2)
// 	if len(s) == 2 {
// 		if len(s[0]) > 0 && s[0][0] != '/' {
// 			s[0] = "/" + s[0]
// 		}
// 		return s[0], s[1]
// 	}
// 	return s[0], ""
// }

// func readPtr(rc io.ReadCloser, ptr string) (io.ReadCloser, error) {
// 	panic("not impl")
// 	// p, err := jsonptr.NewJsonPointer(ptr)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// r, e, err := detectEncoding(rc)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// var b []byte
// 	// if e == yamlEncoding {
// 	// 	b, err = ioutil.ReadAll(r)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	b, err = yamlutil.YAMLToJSON(b)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	r = bytes.NewBuffer(b)
// 	// }
// 	// var c interface{}
// 	// if err = json.NewDecoder(r).Decode(&c); err != nil {
// 	// 	return nil, err
// 	// }
// 	// v, _, err := p.Get(c)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// buf := &bytes.Buffer{}
// 	// enc := json.NewEncoder(buf)
// 	// if err = enc.Encode(v); err != nil {
// 	// 	return nil, err
// 	// }
// 	// return &readercloser{
// 	// 	Reader: buf,
// 	// 	Closer: rc,
// 	// }, nil
// }

// type cache struct {
// 	Params          map[string]*Parameter
// 	Responses       map[string]*Response
// 	Examples        map[string]*Example
// 	Headers         map[string]*Header
// 	RequestBodies   map[string]*RequestBody
// 	Callbacks       map[string]*Callback
// 	Paths           map[string]*Path
// 	SecuritySchemes map[string]*SecurityScheme
// 	Links           map[string]*Link
// 	Schemas         map[string]*Schema
// }

// func newCache() *cache {
// 	return &cache{
// 		Params:          make(map[string]*Parameter),
// 		Responses:       make(map[string]*Response),
// 		Examples:        make(map[string]*Example),
// 		Headers:         make(map[string]*Header),
// 		RequestBodies:   make(map[string]*RequestBody),
// 		Callbacks:       make(map[string]*Callback),
// 		Paths:           make(map[string]*Path),
// 		SecuritySchemes: make(map[string]*SecurityScheme),
// 		Links:           make(map[string]*Link),
// 		Schemas:         make(map[string]*Schema),
// 	}
// }
