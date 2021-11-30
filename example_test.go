package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	assert := require.New(t)
	j := [][]byte{
		[]byte(`{
			"foo": {
			  "summary": "A foo example",
			  "value": {
				"foo": "bar"
			  }
			},
			"bar": {
			  "summary": "A bar example",
			  "value": {
				"bar": "baz"
			  }
			}
		  }`),
		[]byte(`{
			"zip-example": {
			  "$ref": "#/components/examples/zip-example"
			}
		  }`),
		[]byte(`{
			"confirmation-success": {
			  "$ref": "#/components/examples/confirmation-success"
			}
		  }`),
	}
	for _, data := range j {
		var e openapi.Examples
		err := json.Unmarshal(data, &e)
		assert.NoError(err)
		b, err := json.Marshal(e)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b))
	}
}
