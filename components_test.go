package openapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/chanced/openapi"
)

func TestComponentsMarshaling(t *testing.T) {
	f, err := testdata.Open("testdata/schemas/petstore-schema-map-test-1.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fc, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	b := bytes.Buffer{}

	b.Write([]byte(`{"schemas":`))
	b.Write(fc)
	b.Write([]byte(`}`))

	var c openapi.Components
	err = c.UnmarshalJSON(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	cb, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Error(err)
	}
	_ = cb
}
