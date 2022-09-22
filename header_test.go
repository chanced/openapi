package openapi_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	"github.com/chanced/cmpjson"
// 	"github.com/chanced/openapi"
// 	jsonpatch "github.com/evanphx/json-patch/v5"
// 	"github.com/stretchr/testify/require"
// 	yaml "sigs.k8s.io/yaml"
// )

// func TestHeader(t *testing.T) {
// 	assert := require.New(t)

// 	j := []string{
// 		`{
// 		"description": "The number of allowed requests in the current period",
// 		"schema": {
// 		  "type": "integer"
// 		},
// 		"x-header-ext": "value"
// 	  }`,
// 	}
// 	for _, d := range j {
// 		data := []byte(d)
// 		var h openapi.Header
// 		err := json.Unmarshal(data, &h)
// 		assert.NoError(err)

// 		b, err := json.MarshalIndent(h, "", "  ")
// 		assert.NoError(err)
// 		p, err := jsonpatch.CreateMergePatch(data, b)
// 		assert.NoError(err)
// 		assert.True(jsonpatch.Equal(data, b), string(p))

// 		// testing yaml

// 		y, err := yaml.JSONToYAML(data)
// 		assert.NoError(err)
// 		var yo openapi.Header
// 		err = yaml.Unmarshal(y, &yo)
// 		assert.NoError(err)
// 		yb, err := json.MarshalIndent(yo, "", "  ")
// 		assert.NoError(err)
// 		if !jsonpatch.Equal(data, yb) {
// 			fmt.Println(string(data), "\n------------------------\n", string(yb))
// 		}
// 		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

// 	}
// }
