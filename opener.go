package openapi

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
)

type openiniter interface {
	Init(string) error
}

// Openers is a map of URI addresses to Openers.
type Openers map[string]Opener

// Opener is implemented by any value that has an Open method, which accepts a
// path (string) and returns an io.ReadCloser.
//
// The path provided will be in the form of "./path-to-resource", relative to
// the key provided to the Openers map.
//
// Two Openers are provided:
// 	- openapi.FSOpener: opens from an fs.FS
// 	- openapi.HTTPOpener: opens by making HTTP requests
type Opener interface {
	Open(path string) (io.ReadCloser, error)
}

func (oar *OpenAPIResolver) open(pth string) (io.ReadCloser, error) {
	u, o, err := oar.opener(pth)
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

// FSOpener opens files from a filesystem
type FSOpener struct {
	FS fs.FS
}

// Open an OpenAPI object
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

// Init initializes the HTTPOpener. If an error occurs, it will be returned
// uponn a call to Open.
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
