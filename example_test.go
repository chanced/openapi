package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml"
)

func TestExample(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
			"foo": {
			  "summary": "A foo example",
			  "value": {
				"foo": "bar"
			  }
			},
			"bar": {
			  "summary": "A bar example",
			  "value": {
				"bar": "baz"
			  }
			}
		  }`,
		`{
			"zip-example": {
			  "$ref": "#/components/examples/zip-example"
			}
		  }`,
		`{
			"confirmation-success": {
			  "$ref": "#/components/examples/confirmation-success"
			}
		  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var e openapi.Examples
		err := json.Unmarshal(data, &e)
		assert.NoError(err)
		b, err := json.Marshal(e)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b))

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.Examples
		err = yaml.Unmarshal(y, &yo)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yo, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n------------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

	}
}
