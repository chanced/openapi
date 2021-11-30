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

func TestResponse(t *testing.T) {
	assert := require.New(t)

	j := []string{
		`{
		"description": "A complex object array response",
		"content": {
		  "application/json": {
			"schema": {
			  "type": "array",
			  "items": {
				"$ref": "#/components/schemas/VeryComplexType"
			  }
			}
		  }
		}
	  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var v openapi.ResponseObj
		err := json.Unmarshal(data, &v)
		assert.NoError(err)
		b, err := json.Marshal(v)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.ResponseObj
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
