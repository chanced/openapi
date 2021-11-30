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

func TestMediaType(t *testing.T) {
	assert := require.New(t)
	j := []string{`{
		"application/json": {
		  "schema": {
			   "$ref": "#/components/schemas/Pet"
		  },
		  "examples": {
			"cat" : {
			  "summary": "An example of a cat",
			  "value": {
				  "name": "Fluffy",
				  "petType": "Cat",
				  "color": "White",
				  "gender": "male",
				  "breed": "Persian"
				}
			},
			"dog": {
			  "summary": "An example of a dog with a cat's name",
			  "value" :  { 
				"name": "Puma",
				"petType": "Dog",
				"color": "Black",
				"gender": "Female",
				"breed": "Mixed"
			  }
			},
			"frog": {
				"$ref": "#/components/examples/frog-example"
			  }
			}
		}
	  }`,
	}
	for _, d := range j {
		data := []byte(d)
		var mt openapi.Content
		err := json.Unmarshal(data, &mt)
		assert.NoError(err)

		b, err := json.MarshalIndent(mt, "", "  ")
		assert.NoError(err)

		p, err := jsondiff.CompareJSON(data, b)
		if !jsonpatch.Equal(data, b) {
			fmt.Println(string(data))
			fmt.Println(string(b))
		}
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), p.String())

		// testing yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.Content
		err = yaml.Unmarshal(y, &yo)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yo, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n------------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

	}
}
