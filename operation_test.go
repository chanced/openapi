package openapi_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml"
)

func TestOperation(t *testing.T) {
	assert := require.New(t)

	j := []string{
		`{
		"tags": [
		  "pet"
		],
		"summary": "Updates a pet in the store with form data",
		"operationId": "updatePetWithForm",
		"parameters": [
		  {
			"name": "petId",
			"in": "path",
			"description": "ID of pet that needs to be updated",
			"required": true,
			"schema": {
			  "type": "string"
			}
		  }
		],
		"requestBody": {
		  "content": {
			"application/x-www-form-urlencoded": {
			  "schema": {
				"type": "object",
				"properties": {
				  "name": { 
					"description": "Updated name of the pet",
					"type": "string"
				  },
				  "status": {
					"description": "Updated status of the pet",
					"type": "string"
				  }
				},
				"required": ["status"] 
			  }
			}
		  }
		},
		"responses": {
		  "200": {
			"description": "Pet updated.",
			"content": {
			  "application/json": {},
			  "application/xml": {}
			}
		  },
		  "405": {
			"description": "Method Not Allowed",
			"content": {
			  "application/json": {},
			  "application/xml": {}
			}
		  }
		},
		"security": [
		  {
			"petstore_auth": [
			  "write:pets",
			  "read:pets"
			]
		  }
		]
	  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var o openapi.Operation
		err := json.Unmarshal(data, &o)
		assert.NoError(err)
		b, err := json.Marshal(o)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))

		// testing yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.Operation
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

func TestExtensionSorting(t *testing.T) {
	assert := require.New(t)
	exp := `{"x-key1":1,"x-key2":2}`
	for n := 0; n < 100; n++ {
		op := new(openapi.Operation)
		op.Extensions = make(openapi.Extensions)
		op.Extensions.SetEncodedExtension("key1", []byte("1"))
		op.Extensions.SetEncodedExtension("key2", []byte("2"))

		marshaled, _ := json.Marshal(op)

		assert.Equal(string(marshaled), exp, strconv.Itoa(n))
	}
}
