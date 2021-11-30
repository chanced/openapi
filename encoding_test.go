package openapi_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	"github.com/chanced/openapi/yamlutil"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	assert := require.New(t)

	y := bytes.ReplaceAll([]byte(`
historyMetadata:
	# require XML Content-Type in utf-8 encoding
	contentType: application/xml; charset=utf-8
profileImage:
	# only accept png/jpeg
	contentType: image/png, image/jpeg
	headers:
		X-Rate-Limit-Limit:
			description: The number of allowed requests in the current period
			schema:
				type: integer
`), []byte("\t"), []byte("    "))
	d, err := yamlutil.YAMLToJSON(y)
	assert.NoError(err)
	var e openapi.Encodings
	err = json.Unmarshal(d, &e)
	assert.NoError(err)
	b, err := json.MarshalIndent(e, "", "  ")
	data := []byte(`{
		"historyMetadata": {
		  "contentType": "application/xml; charset=utf-8"
		},
		"profileImage": {
		  "contentType": "image/png, image/jpeg",
		  "headers": {
			"X-Rate-Limit-Limit": {
			  "description": "The number of allowed requests in the current period",
			  "schema": {
				"type": "integer"
			  }
			}
		  }
		}
	  }`)
	assert.True(jsonpatch.Equal(data, b))
	assert.NoError(err)
}
