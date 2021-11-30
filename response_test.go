package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	assert := require.New(t)

	j := []byte(`{
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
	  }`)
	var v openapi.ResponseObj

	err := json.Unmarshal(j, &v)
	assert.NoError(err)

	b, err := json.Marshal(v)
	assert.NoError(err)
	assert.True(jsonpatch.Equal(j, b), cmpjson.Diff(j, b))

}
