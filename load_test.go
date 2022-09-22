package openapi_test

import (
	"context"
	"embed"
	"fmt"
	"io"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/chanced/openapi"
	"github.com/chanced/transcode"
	"github.com/chanced/uri"
)

//go:embed testdata
var testdata embed.FS

func TestLoadRefComponent(t *testing.T) {
	f, err := testdata.Open("testdata/documents/comprefs.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	doc, err := openapi.Load(ctx, "testdata/documents/comprefs.yaml", NoopValidator{}, func(ctx context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error) {
		b, err := io.ReadAll(f)
		if err != nil {
			return 0, nil, err
		}
		return openapi.KindDocument, b, nil
	})
	if err != nil {
		t.Error(err)
	}
	if doc == nil {
		t.Errorf("failed to load document")
	}
	// litter.Dump(doc)
	if doc.Components.Responses.Get("Referenced").Object.Description != "/components/responses/Referenced" {
		t.Errorf("expected %q got %q", "/components/responses/Referenced", doc.Components.Responses.Get("Referenced").Object.Description)
	}
	refpath := doc.Paths.Get("/ref")
	if refpath.Post.Responses.Get("200").Object.Description != "/components/responses/Referenced" {
		t.Errorf("expected %q got %q", "/components/responses/Referenced", doc.Paths.Get("/refs").Post.Responses.Get("200").Object.Description)
	}
	rb := doc.Components.RequestBodies.Get("Referenced")
	if rb.Object.Description != "/components/requestBodies/Referenced" {
		t.Errorf("expected requestBody to have description of %q, got %q", "/components/requestBodies/Referenced", rb.Object.Description)
	}
	rbr := refpath.Post.RequestBody.Object
	if rbr.Description != rb.Object.Description {
		t.Errorf("expected requestBody to have description of %q, got %q", rb.Object.Description, rbr.Description)
	}
}

func TestLoad(t *testing.T) {
	f, err := testdata.Open("testdata/documents/petstore.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	doc, err := openapi.Load(ctx, "testdata/documents/petstore.yaml", NoopValidator{}, func(ctx context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error) {
		b, err := io.ReadAll(f)
		// fmt.Println(string(b))
		if err != nil {
			return 0, nil, err
		}
		return openapi.KindDocument, b, nil
	})
	if err != nil {
		t.Error(err)
	}
	if doc == nil {
		t.Errorf("failed to load document")
	}
	// litter.Dump(doc)
}

func TestDynamicRefs(t *testing.T) {
	f, err := testdata.Open("testdata/documents/dynamic-refs.yaml")
	if err != nil {
		t.Fatal(err)
	}
	listOfT, err := testdata.Open("testdata/schemas/list-of-t.json")
	if err != nil {
		t.Fatal(err)
	}
	listOfStrings, err := testdata.Open("testdata/schemas/list-of-strings.json")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	doc, err := openapi.Load(ctx, "testdata/documents/dynamic-refs.yaml", NoopValidator{}, func(ctx context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error) {
		switch uri.String() {
		case "https://json-schema.blog/list-of-t":
			d, err := io.ReadAll(listOfT)
			return openapi.KindSchema, d, err
		case "https://json-schema.blog/list-of-strings":
			d, err := io.ReadAll(listOfStrings)
			return openapi.KindSchema, d, err
		case "testdata/documents/dynamic-refs.yaml":
			d, err := io.ReadAll(f)
			return openapi.KindDocument, d, err
		default:
			return 0, nil, fmt.Errorf("uknown uri %q", uri)
		}
	})
	if err != nil {
		t.Error(err)
	}

	// litter.Dump(doc)
	los := doc.Components.Schemas.Get("ListOfStrings")
	if los.Ref.Resolved.Ref.Resolved == nil {
		t.Error("expected /list-of-t to be resolved")
	}
}

type NoopValidator struct{}

func (NoopValidator) Validate(data []byte, resource uri.URI, kind openapi.Kind, openapi semver.Version, jsonschema uri.URI) error {
	return nil
}

func TestTryGetOpenAPIVersion(t *testing.T) {
	f, err := testdata.Open("testdata/documents/petstore.yaml")
	if err != nil {
		t.Fatal(err)
	}
	d, _ := io.ReadAll(f)
	d, err = transcode.JSONFromYAML(d)
	if err != nil {
		t.Errorf("failed to transcode data")
	}
	if len(d) == 0 {
		t.Fatal("file was empty")
	}
	vstr, ok := openapi.TryGetOpenAPIVersion(d)
	if !ok {
		t.Error("failed to get openapi")
	}
	if vstr != "3.1.0" {
		t.Errorf("expected 3.1.0 got %q", vstr)
	}
}

func (NoopValidator) ValidateDocument(document *openapi.Document) error { return nil }

var _ openapi.Validator = (*NoopValidator)(nil)
