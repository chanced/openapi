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

func TestEncoding(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
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
		  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var e openapi.Encodings
		err := json.Unmarshal(data, &e)
		assert.NoError(err)
		b, err := json.MarshalIndent(e, "", "  ")
		assert.NoError(err)
		d, err := jsondiff.CompareJSON(data, b)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), d.String())

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.Encodings
		err = yaml.Unmarshal(y, &yo)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yo, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n------------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

		assert.True(jsonpatch.Equal(data, b))
		assert.NoError(err)
	}
}
