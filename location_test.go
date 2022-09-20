package openapi_test

import (
	"testing"

	"github.com/chanced/openapi"
	"github.com/chanced/uri"
)

func TestLocationAppend(t *testing.T) {
	u, _ := uri.Parse("https://example.org/schema/demo")

	loc, err := openapi.NewLocation(*u)
	if err != nil {
		t.Fatal(err)
	}
	loc = loc.AppendLocation("foo")
	expected := "https://example.org/schema/demo#/foo"
	if loc.String() != expected {
		t.Errorf("expected %q, got %s", expected, loc.String())
	}

	loc = loc.AppendLocation("bar")
	expected = "https://example.org/schema/demo#/foo/bar"
	if loc.String() != expected {
		t.Errorf("expected %q, got %s", expected, loc.String())
	}
	u, err = uri.Parse("example.json")
	if err != nil {
		t.Fatal(err)
	}
	loc, err = openapi.NewLocation(*u)
	if err != nil {
		t.Fatal(err)
	}
	loc = loc.AppendLocation("foo")
	expected = "example.json#/foo"
	if loc.String() != expected {
		t.Errorf("expected %q, got %s", expected, loc.String())
	}
}
