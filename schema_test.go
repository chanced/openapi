package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	"github.com/chanced/uri"
)

func TestSchema(t *testing.T) {
	bs := []byte(`true`)

	var s openapi.Schema
	err := json.Unmarshal(bs, &s)
	if err != nil {
		t.Error(err)
	}
	sb, err := s.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	if string(sb) != "true" {
		t.Errorf("expected %q, got %q", "true", string(sb))
	}
	var s2 openapi.Schema
	bs = []byte(`{"keyword": "value"}`)
	err = json.Unmarshal(bs, &s2)
	if err != nil {
		t.Error(err)
	}
	if string(s2.Keywords["keyword"]) != `"value"` {
		t.Errorf("expected %q, got %q", "value", s2.Keywords["keyword"])
	}
	br, err := s2.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	_ = br
	// fmt.Println(string(br))
}

func TestClone(t *testing.T) {
	s := openapi.Schema{
		If: &openapi.Schema{
			Format: "format",
		},
		Schema: uri.MustParse("https://json-schema.org/draft/2019-09/schema"),
		ID:     uri.MustParse("http://example.com/schema"),
		Ref: &openapi.SchemaRef{
			Ref: uri.MustParse("http://example.com/schema2"),
			Resolved: &openapi.Schema{
				ID: uri.MustParse("http://example.com/resolved"),
			},
		},
	}
	s2 := s.Clone()
	if s2.If.Format != "format" {
		t.Errorf("expected %q, got %q", "format", s2.If.Format)
	}
	if s2.Ref.Resolved == nil {
		t.Error("expected resolved schema, got nil")
	}
}
