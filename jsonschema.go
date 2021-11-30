package openapi

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

//go:embed schema
var schemaDir embed.FS
var schemas = map[string]io.ReadCloser{}

func init() {
	if err := loadSchemas(); err != nil {
		log.Fatal(err)
	}
	jsonschema.Loaders["map"] = loadSchema
}

func loadSchema(url string) (io.ReadCloser, error) {
	s, ok := schemas[url]
	if !ok {
		return nil, fmt.Errorf("schema not found: %s", url)
	}
	return s, nil
}

func loadSchemas() error {
	return fs.WalkDir(schemaDir, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || filepath.Ext(d.Name()) != "json" {
			return nil
		}
		f, err := schemaDir.Open(path)
		if err != nil {
			return err
		}
		j, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		g := gjson.GetBytes(j, "$id")
		if !g.Exists() || len(g.String()) == 0 {
			return errors.New("schema is missing $id")
		}
		schemas[g.String()] = ioutil.NopCloser(bytes.NewReader(j))
		return nil
	})
}
