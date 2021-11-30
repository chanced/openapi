package openapi

import (
	"encoding/json"
	"testing"

	"github.com/chanced/cmpjson"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestOperation(t *testing.T) {
	assert := require.New(t)

	j := [][]byte{
		[]byte(`{
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
	  }`),
	}
	for _, data := range j {
		var o Operation
		err := json.Unmarshal(data, &o)
		assert.NoError(err)
		b, err := json.Marshal(o)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))
	}
}
