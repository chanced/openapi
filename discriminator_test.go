package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	"github.com/wI2L/jsondiff"
)

func TestDiscriminator(t *testing.T) {
	assert := require.New(t)
	j := []byte(`{
		"propertyName": "petType",
		"mapping": {
			"dog": "#/components/schemas/Dog",
			"monster": "https://gigantic-server.com/schemas/Monster/schema.json"
		},
		"x-ext":  "ext val",
		"x-ext2": 2
	}`)

	var dis openapi.Discriminator
	err := json.Unmarshal(j, &dis)
	assert.NoError(err)
	b, err := json.MarshalIndent(dis, "", "  ")
	assert.NoError(err)
	diff, err := jsondiff.CompareJSON(j, b)
	assert.NoError(err)
	assert.True(jsonpatch.Equal(j, b), diff.String())
	div := &openapi.Discriminator{}
	div.PropertyName = "prop"

	_, err = json.MarshalIndent(div, "", "  ")
	assert.NoError(err)

}
