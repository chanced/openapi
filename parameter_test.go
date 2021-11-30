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

func TestParameter(t *testing.T) {
	assert := require.New(t)

	j := []string{
		`{
			"name": "token",
			"in": "header",
			"description": "token to be passed as a header",
			"required": true,
			"schema": {
			  "type": "array",
			  "items": {
				"type": "integer",
				"format": "int64"
			  }
			},
			"style": "simple"
		  }`,
		`{
			"name": "username",
			"in": "path",
			"description": "username to fetch",
			"required": true,
			"schema": {
			  "type": "string"
			}
		  }`,
		`{
			"name": "id",
			"in": "query",
			"description": "ID of the object to fetch",
			"required": false,
			"schema": {
			  "type": "array",
			  "items": {
				"type": "string"
			  }
			},
			"style": "form",
			"explode": true
		  }`,
		`{
			"in": "query",
			"name": "freeForm",
			"schema": {
			  "type": "object",
			  "additionalProperties": {
				"type": "integer"
			  }
			},
			"style": "form"
		  }`,
		`{
			"in": "query",
			"name": "coordinates",
			"content": {
			  "application/json": {
				"schema": {
				  "type": "object",
				  "required": [
					"lat",
					"long"
				  ],
				  "properties": {
					"lat": {
					  "type": "number"
					},
					"long": {
					  "type": "number"
					}
				  }
				}
			  }
			}
		  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var p openapi.ParameterObj
		err := json.Unmarshal(data, &p)
		assert.NoError(err)
		b, err := json.MarshalIndent(p, "", "  ")

		assert.NoError(err)
		if !jsonpatch.Equal(data, b) {
			fmt.Println(string(b))
		}

		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.ParameterObj
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
