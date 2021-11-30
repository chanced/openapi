package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestLink(t *testing.T) {
	assert := require.New(t)
	j := [][]byte{
		[]byte(`{
			"address": {
			  "operationId": "getUserAddress",
			  "parameters": {
				"userId": "$request.path.id"
			  }
			}
		  }`),
		[]byte(`{
			"UserRepositories": {
			  "operationRef": "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get",
			  "parameters": {
				"username": "$response.body#/username"
			  }
			}
		  }`),
	}

	for _, data := range j {
		var ll openapi.Links
		err := json.Unmarshal(data, &ll)
		assert.NoError(err)
		b, err := json.Marshal(ll)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b))
	}
}
