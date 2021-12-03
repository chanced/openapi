package openapi_test

import (
	"testing"

	"github.com/chanced/openapi"
	"github.com/stretchr/testify/require"
)

func TestReference(t *testing.T) {
	assert := require.New(t)
	r := "http://example.com/test.json"
	var v openapi.Callback = &openapi.Reference{
		Ref: r,
	}
	ran := false
	v.ResolveCallback(func(ref string) (*openapi.CallbackObj, error) {
		ran = true
		assert.Equal(r, ref)
		return &openapi.CallbackObj{}, nil
	})
	assert.True(ran)
}
