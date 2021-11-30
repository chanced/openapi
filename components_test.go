package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	"github.com/wI2L/jsondiff"
	yaml "sigs.k8s.io/yaml"
)

func TestComponents(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
		"schemas": {
		  "GeneralError": {
			"type": "object",
			"properties": {
			  "code": {
				"type": "integer",
				"format": "int32"
			  },
			  "message": {
				"type": "string"
			  }
			}
		  },
		  "Category": {
			"type": "object",
			"properties": {
			  "id": {
				"type": "integer",
				"format": "int64"
			  },
			  "name": {
				"type": "string"
			  }
			}
		  },
		  "Tag": {
			"type": "object",
			"properties": {
			  "id": {
				"type": "integer",
				"format": "int64"
			  },
			  "name": {
				"type": "string"
			  }
			}
		  }
		},
		"parameters": {
		  "skipParam": {
			"name": "skip",
			"in": "query",
			"description": "number of items to skip",
			"required": true,
			"schema": {
			  "type": "integer",
			  "format": "int32"
			}
		  },
		  "limitParam": {
			"name": "limit",
			"in": "query",
			"description": "max records to return",
			"required": true,
			"schema" : {
			  "type": "integer",
			  "format": "int32"
			}
		  }
		},
		"responses": {
		  "NotFound": {
			"description": "Entity not found."
		  },
		  "IllegalInput": {
			"description": "Illegal input for operation."
		  },
		  "GeneralError": {
			"description": "General Error",
			"content": {
			  "application/json": {
				"schema": {
				  "$ref": "#/components/schemas/GeneralError"
				}
			  }
			}
		  }
		},
		"securitySchemes": {
		  "api_key": {
			"type": "apiKey",
			"name": "api_key",
			"in": "header"
		  },
		  "petstore_auth": {
			"type": "oauth2",
			"flows": {
			  "implicit": {
				"authorizationUrl": "https://example.org/api/oauth/dialog",
				"scopes": {
				  "write:pets": "modify pets in your account",
				  "read:pets": "read your pets"
				}
			  }
			}
		  }
		}
	  }`,
	}

	for _, d := range j {
		data := []byte(d)
		var c openapi.Components
		err := json.Unmarshal(data, &c)
		assert.NoError(err)
		b, err := json.MarshalIndent(c, "", "  ")
		assert.NoError(err)
		d, err := jsondiff.CompareJSON(data, b)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), d.String())

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yc openapi.Components
		err = yaml.Unmarshal(y, &yc)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yc, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n-----------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

	}
}
