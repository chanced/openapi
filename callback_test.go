package openapi_test

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/chanced/cmpjson"
// 	"github.com/chanced/openapi"
// 	jsonpatch "github.com/evanphx/json-patch/v5"
// 	"github.com/stretchr/testify/require"
// )

// func TestCallback(t *testing.T) {
// 	assert := require.New(t)
// 	cbj := [][]byte{
// 		[]byte(`{
// 		  "{$request.query.queryUrl}": {
// 			"post": {
// 			  "requestBody": {
// 				"description": "Callback payload",
// 				"content": {
// 				  "application/json": {
// 					"schema": {
// 					  "$ref": "#/components/schemas/SomePayload"
// 					}
// 				  }
// 				}
// 			  },
// 			  "responses": {
// 				"200": {
// 				  "description": "callback successfully processed"
// 				}
// 			  }
// 			}
// 		  },
// 		  "x-test": "value"
// 	  }`),
// 	}
// 	for _, data := range cbj {
// 		var v openapi.Callback
// 		err := json.Unmarshal(data, &v)
// 		assert.NoError(err)
// 		b, err := json.MarshalIndent(&v, "", "  ")
// 		assert.NoError(err)
// 		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))
// 	}

// 	cbjs := [][]byte{
// 		[]byte(`{
// 		"myCallback": {
// 		  "{$request.query.queryUrl}": {
// 			"post": {
// 			  "requestBody": {
// 				"description": "Callback payload",
// 				"content": {
// 				  "application/json": {
// 					"schema": {
// 					  "$ref": "#/components/schemas/SomePayload"
// 					}
// 				  }
// 				}
// 			  },
// 			  "responses": {
// 				"200": {
// 				  "description": "callback successfully processed"
// 				}
// 			  }
// 			}
// 		  }
// 		}
// 	  }`),
// 	}
// 	for _, data := range cbjs {
// 		var v openapi.CallbackMap
// 		err := json.Unmarshal(data, &v)
// 		assert.NoError(err)
// 		b, err := json.MarshalIndent(&v, "", "  ")
// 		assert.NoError(err)
// 		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))
// 	}
// }
