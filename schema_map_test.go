package openapi_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/chanced/openapi"
)

func TestSchemaMapMarshaling(t *testing.T) {
	f, err := testdata.Open("testdata/schemas/petstore-schema-map-test-1.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fc, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	var sm openapi.SchemaMap
	err = json.Unmarshal(fc, &sm)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
