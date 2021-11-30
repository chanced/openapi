package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	"github.com/wI2L/jsondiff"
)

func TestMediaType(t *testing.T) {
	assert := require.New(t)
	j := []byte(`{
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
	  }`)
	var mt openapi.Content
	err := json.Unmarshal(j, &mt)
	assert.NoError(err)

	b, err := json.MarshalIndent(mt, "", "  ")
	assert.NoError(err)

	p, err := jsondiff.CompareJSON(j, b)
	if !jsonpatch.Equal(j, b) {
		fmt.Println(string(j))
		fmt.Println(string(b))
	}
	assert.NoError(err)
	assert.True(jsonpatch.Equal(j, b), p.String())
}
