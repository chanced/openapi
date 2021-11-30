package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	"github.com/wI2L/jsondiff"
)

func TestContact(t *testing.T) {
	assert := require.New(t)
	j := [][]byte{
		[]byte(`{
			"name": "API Support",
			"url": "https://www.example.com/support",
			"email": "support@example.com"
		}`),
	}

	for _, data := range j {
		var c openapi.Contact
		err := json.Unmarshal(data, &c)
		assert.NoError(err)
		b, err := json.MarshalIndent(c, "", "  ")
		assert.NoError(err)
		p, err := jsondiff.CompareJSON(data, b)
		assert.NoError(err)
		assert.True(jsonpatch.Equal(data, b), p)
	}
}
