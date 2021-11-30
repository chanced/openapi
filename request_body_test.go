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

func TestRequestBody(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
		"description": "user to add to the system",
		"content": {
		  "application/json": {
			"schema": {
			  "$ref": "#/components/schemas/User"
			},
			"examples": {
				"user" : {
				  "summary": "User Example", 
				  "externalValue": "https://foo.bar/examples/user-example.json"
				} 
			  }
		  },
		  "application/xml": {
			"schema": {
			  "$ref": "#/components/schemas/User"
			},
			"examples": {
				"user" : {
				  "summary": "User example in XML",
				  "externalValue": "https://foo.bar/examples/user-example.xml"
				}
			  }
		  },
		  "text/plain": {
			"examples": {
			  "user" : {
				  "summary": "User example in Plain text",
				  "externalValue": "https://foo.bar/examples/user-example.txt" 
			  }
			} 
		  },
		  "*/*": {
			"examples": {
			  "user" : {
				  "summary": "User example in other format",
				  "externalValue": "https://foo.bar/examples/user-example.whatever"
			  }
			}
		  }
		}
	  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var rb openapi.RequestBodyObj
		err := json.Unmarshal(data, &rb)
		assert.NoError(err)
		b, err := json.Marshal(rb)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.RequestBodyObj
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
