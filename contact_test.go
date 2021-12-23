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

func TestContact(t *testing.T) {
	assert := require.New(t)
	j := []string{
		`{
			"name": "API Support",
			"url": "https://www.example.com/support",
			"email": "support@example.com"
		}`,
	}

	for _, d := range j {
		data := []byte(d)
		var c openapi.Contact
		err := json.Unmarshal(data, &c)
		assert.NoError(err)
		b, err := json.MarshalIndent(c, "", "  ")
		assert.NoError(err)
		p, err := jsondiff.CompareJSON(data, b)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), p)

		// testing yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yc openapi.Contact
		err = yaml.Unmarshal(y, &yc)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yc, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n------------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

	}
}
