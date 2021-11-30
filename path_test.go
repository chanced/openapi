package openapi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	assert := require.New(t)

	j := []byte(`{
		"/users/{id}": {
		  "parameters": [
			{
			  "name": "id",
			  "in": "path",
			  "required": true,
			  "description": "the user identifier, as userId",
			  "schema": {
				"type": "string"
			  }
			}
		  ],
		  "get": {
			"responses": {
			  "200": {
				"description": "the user being returned",
				"content": {
				  "application/json": {
					"schema": {
					  "type": "object",
					  "properties": {
						"uuid": {
						  "type": "string",
						  "format": "uuid"
						}
					  }
					}
				  }
				},
				"links": {
				  "address": {
					"operationId": "getUserAddress",
					"parameters": {
					  "userId": "$request.path.id"
					}
				  }
				}
			  }
			}
		  }
		},
		"/users/{userid}/address": {
		  "parameters": [
			{
			  "name": "userid",
			  "in": "path",
			  "required": true,
			  "description": "the user identifier, as userId",
			  "schema": {
				"type": "string"
			  }
			}
		  ],
		  "get": {
			"operationId": "getUserAddress",
			"responses": {
			  "200": {
				"description": "the user's address"
			  }
			}
		  }
		}
	  }`)

	var paths openapi.Paths
	err := json.Unmarshal(j, &paths)
	var te *json.UnmarshalTypeError
	if errors.As(err, &te) {
		fmt.Println(te.Field)
		fmt.Println(te.Value)
		fmt.Println(te.Struct)
	}
	assert.NoError(err)

	b, err := json.MarshalIndent(paths, "", "  ")
	assert.NoError(err)
	// patch, err := jsonpatch.CreateMergePatch(j, d)
	assert.NoError(err)
	if !jsonpatch.Equal(j, b) {
		fmt.Println(string(b))
	}
	assert.True(jsonpatch.Equal(j, b), cmpjson.Diff(j, b))
}
