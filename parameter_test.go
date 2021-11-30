package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestParameter(t *testing.T) {
	assert := require.New(t)

	j := [][]byte{
		[]byte(`{
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
		  }`),
		[]byte(`{
			"name": "username",
			"in": "path",
			"description": "username to fetch",
			"required": true,
			"schema": {
			  "type": "string"
			}
		  }`),
		[]byte(`{
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
		  }`),
		[]byte(`{
			"in": "query",
			"name": "freeForm",
			"schema": {
			  "type": "object",
			  "additionalProperties": {
				"type": "integer"
			  }
			},
			"style": "form"
		  }`),
		[]byte(`{
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
		  }`),
	}
	for _, data := range j {
		var p openapi.ParameterObj
		err := json.Unmarshal(data, &p)
		assert.NoError(err)
		b, err := json.MarshalIndent(p, "", "  ")

		assert.NoError(err)
		if !jsonpatch.Equal(data, b) {
			fmt.Println(string(b))
		}

		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))
	}
}
