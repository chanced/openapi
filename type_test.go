package openapi

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTypesUnmarshalYAML(t *testing.T) {
	single := `string`
	multi := `
- string
- number
`
	var st Types
	err := yaml.Unmarshal([]byte(single), &st)
	if err != nil {
		t.Error(err)
	}
	if len(st) != 1 {
		t.Errorf("expected 1 type, got %d", len(st))
	}
	if st[0] != TypeString {
		t.Errorf("expected %q, got %q", TypeString, st[0])
	}
	var mt Types
	err = yaml.Unmarshal([]byte(multi), &mt)
	if err != nil {
		t.Error(err)
	}
	if len(mt) != 2 {
		t.Errorf("expected 2 types, got %d", len(mt))
	}
	if mt[0] != TypeString {
		t.Errorf("expected %q, got %q", TypeString, mt[0])
	}
	if mt[1] != TypeNumber {
		t.Errorf("expected %q, got %q", TypeNumber, mt[1])
	}
}

func TestTypesMarshalYAML(t *testing.T) {
	st := Types{TypeString}

	b, err := yaml.Marshal(st)
	if err != nil {
		t.Error(err)
	}
	if string(b) != "string\n" {
		t.Errorf("expected %q, got %q", "string", string(b))
	}

	mt := Types{TypeString, TypeNumber}
	b, err = yaml.Marshal(mt)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
	if string(b) != "- string\n- number\n" {
		t.Errorf("expected %q, got %q", "- string\n- number", string(b))
	}
}
