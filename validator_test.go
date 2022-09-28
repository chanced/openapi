package openapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/chanced/openapi"
	"github.com/chanced/uri"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

func TestValidation(t *testing.T) {
	ctx := context.Background()

	c, err := openapi.SetupCompiler(jsonschema.NewCompiler()) // adding schema files
	if err != nil {
		t.Fatal(err)
	}
	v, err := openapi.NewValidator(c)
	if err != nil {
		t.Fatal(err)
	}

	// you can Load either JSON or YAML
	err = fs.WalkDir(testdata, "testdata/documents/validation/pass", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		t.Run(strings.TrimPrefix(p, "testdata/documents/validation/"), func(t *testing.T) {
			f, err := testdata.Open(p)
			if err != nil {
				t.Fatal(err)
			}
			fn := func(_ context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error) {
				d, err := io.ReadAll(f)
				if err != nil {
					return 0, nil, err
				}
				return openapi.KindDocument, d, nil
			}
			doc, err := openapi.Load(ctx, p, v, fn)
			if err != nil {
				t.Errorf("failed to load document: %v", err)
			}

			err = v.ValidateDocument(doc)
			if err != nil {
				t.Errorf("failed to validate document: %v", err)
			}
		})
		return nil
	})
	if err != nil {
		t.Errorf("expected document to be valid, received: %v", err)
	}
	// you can Load either JSON or YAML
	err = fs.WalkDir(testdata, "testdata/documents/validation/fail", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		t.Run(strings.TrimPrefix(p, "testdata/documents/validation/"), func(t *testing.T) {
			f, err := testdata.Open(p)
			if err != nil {
				t.Fatal(err)
			}
			fn := func(_ context.Context, uri uri.URI, kind openapi.Kind) (openapi.Kind, []byte, error) {
				d, err := io.ReadAll(f)
				if err != nil {
					return 0, nil, err
				}
				return openapi.KindDocument, d, nil
			}
			_, err = openapi.Load(ctx, p, v, fn)
			if err == nil {
				t.Errorf("expected document to be invalid, received: %v", err)
			}

			var ve *openapi.ValidationError
			if !errors.As(err, &ve) {
				t.Errorf("expected document to be invalid, received: %v", err)
			}
		})
		return nil
	})
	if err != nil {
		t.Errorf("expected document to be valid, received: %v", err)
	}
}

// issue #19
// https://github.com/chanced/openapi/issues/19
func TestValidator_Schema(t *testing.T) {
	// ctx := context.Background()

	c, err := openapi.SetupCompiler(jsonschema.NewCompiler()) // adding schema files
	if err != nil {
		t.Fatal(err)
	}
	v, err := openapi.NewValidator(c)
	if err != nil {
		t.Fatal(err)
	}
	f, err := testdata.Open("testdata/schemas/string-map.yaml")
	if err != nil {
		t.Fatal(err)
	}
	d, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	var s openapi.Schema
	err = yaml.Unmarshal(d, &s)
	if err != nil {
		t.Fatal(err)
	}
	d, err = json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	v.Validate(d, *uri.MustParse("testdata/schemas/string-map.yaml"), openapi.KindSchema, openapi.Version3_1, openapi.JSONSchemaDialect202012)
}
