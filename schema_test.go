package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/openapi"
)

func TestSchema(t *testing.T) {
	bs := []byte(`true`)

	var s openapi.Schema
	err := json.Unmarshal(bs, &s)
	if err != nil {
		t.Error(err)
	}
	if s.Always == nil || !*s.Always {
		t.Error("expected Always to be true")
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
	fmt.Println(string(br))
}
