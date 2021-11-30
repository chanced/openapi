package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	"github.com/stretchr/testify/require"
)

func TestRequestBody(t *testing.T) {
	assert := require.New(t)
	j := []byte(`{
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
	  }`)

	var rb openapi.RequestBodyObj
	err := json.Unmarshal(j, &rb)
	assert.NoError(err)
}
