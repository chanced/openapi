package openapi

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	"github.com/chanced/cmpjson"
// 	"github.com/chanced/openapi"
// 	jsonpatch "github.com/evanphx/json-patch/v5"
// 	"github.com/stretchr/testify/require"
// 	yaml "sigs.k8s.io/yaml"
// )

// func TestIssue5(t *testing.T) {
// 	assert := require.New(t)
// 	data := `{
// 	  "openapi": "3.1.0",
// 	  "info": {
// 		"title": "",
// 		"version": "",
// 		"description": "Test file for loading pre-existing OAS"
// 	  },
// 	  "paths": {
// 		"/catalogue/{id}": {
// 		  "parameters": [
// 			{
// 			  "name": "id",
// 			  "in": "path",
// 			  "required": true,
// 			  "style": "simple",
// 			  "schema": {
// 				"type": "string"
// 			  },
// 			  "examples": {
// 				"an example": {
// 				  "value": "someval"
// 				}
// 			  }
// 			}
// 		  ]
// 		},
// 		"/catalogue/{id}/details": {
// 		  "parameters": [
// 			{
// 			  "name": "id",
// 			  "in": "path",
// 			  "style": "simple",
// 			  "required": true,
// 			  "schema": {
// 				"type": "string"
// 			  },
// 			  "example": "some-uuid-maybe"
// 			}
// 		  ]
// 		}
// 	  }
// 	}`

// 	var oas openapi.OpenAPI
// 	err := json.Unmarshal([]byte(data), &oas)
// 	assert.NoError(err)
// 	pi := oas.Paths.Items["/catalogue/{id}"]
// 	assert.NotNil(pi)
// 	assert.NotNil(pi.Parameters)
// 	assert.Len(*pi.Parameters, 1)
// 	params := *pi.Parameters
// 	param := params[0]
// 	paramobj := param.Object
// 	assert.Contains(paramobj.Examples, "an example")
// 	ex, ok := paramobj.Examples.Get("an example")
// 	assert.True(ok)
// 	assert.Equal(json.RawMessage(`"someval"`), ex.Object.Value)
// }

// func TestExample(t *testing.T) {
// 	assert := require.New(t)
// 	j := []string{
// 		`{
// 			"foo": {
// 			  "summary": "A foo example",
// 			  "value": {
// 				"foo": "bar"
// 			  }
// 			},
// 			"bar": {
// 			  "summary": "A bar example",
// 			  "value": {
// 				"bar": "baz"
// 			  }
// 			}
// 		  }`,
// 		`{
// 			"zip-example": {
// 			  "$ref": "#/components/examples/zip-example"
// 			}
// 		  }`,
// 		`{
// 			"confirmation-success": {
// 			  "$ref": "#/components/examples/confirmation-success"
// 			}
// 		  }`,
// 	}
// 	for _, d := range j {
// 		data := []byte(d)
// 		var e openapi.ExampleMap
// 		err := json.Unmarshal(data, &e)
// 		assert.NoError(err)
// 		b, err := json.Marshal(e)
// 		assert.NoError(err)
// 		assert.True(jsonpatch.Equal(data, b))

// 		// testing yaml

// 		y, err := yaml.JSONToYAML(data)
// 		assert.NoError(err)
// 		var yo openapi.ExampleMap
// 		err = yaml.Unmarshal(y, &yo)
// 		assert.NoError(err)
// 		yb, err := json.MarshalIndent(yo, "", "  ")
// 		assert.NoError(err)
// 		if !jsonpatch.Equal(data, yb) {
// 			fmt.Println(string(data), "\n------------------------\n", string(yb))
// 		}
// 		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

// 	}
// }
