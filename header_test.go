package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	assert := require.New(t)

	j := []byte(`{
		"description": "The number of allowed requests in the current period",
		"schema": {
		  "type": "integer"
		}, 
		"x-header-ext": "value"
	  }`)
	var h openapi.HeaderObj
	err := json.Unmarshal(j, &h)
	assert.NoError(err)

	b, err := json.MarshalIndent(h, "", "  ")
	assert.NoError(err)
	p, err := jsonpatch.CreateMergePatch(j, b)
	assert.NoError(err)
	assert.True(jsonpatch.Equal(j, b), string(p))
}
