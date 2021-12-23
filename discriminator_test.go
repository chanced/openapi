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

func TestDiscriminator(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
		"propertyName": "petType",
		"mapping": {
			"dog": "#/components/schemas/Dog",
			"monster": "https://gigantic-server.com/schemas/Monster/schema.json"
		},
		"x-ext":  "ext val",
		"x-ext2": 2
	}`,
	}
	for _, d := range j {
		data := []byte(d)
		var dis openapi.Discriminator
		err := json.Unmarshal(data, &dis)
		assert.NoError(err)
		b, err := json.MarshalIndent(dis, "", "  ")
		assert.NoError(err)
		diff, err := jsondiff.CompareJSON(data, b)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), diff.String())
		div := &openapi.Discriminator{}
		div.PropertyName = "prop"

		_, err = json.MarshalIndent(div, "", "  ")
		assert.NoError(err)

		// testing yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.Discriminator
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
