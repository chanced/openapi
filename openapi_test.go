package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"

	// "gopkg.in/yaml.v2"
	yaml "sigs.k8s.io/yaml"
)

func TestOpenAPI(t *testing.T) {
	pass := []string{
		`{
      "openapi": "3.1.0",
      "info": {
        "summary": "My API's summary",
        "title": "My API",
        "version": "1.0.0",
        "license": {
          "name": "Apache 2.0",
          "identifier": "Apache-2.0"
        }
      },
      "jsonSchemaDialect": "https://spec.openapis.org/oas/3.1/dialect/base",
      "paths": {
        "/": {
          "get": {
            "parameters": []
          }
        },
        "/{pathTest}": {}
      },
      "webhooks": {
        "myWebhook": {
          "$ref": "#/components/pathItems/myPathItem",
          "description": "Overriding description"
        }
      },
      "components": {
        "securitySchemes": {
          "mtls": {
            "type": "mutualTLS"
          }
        },
        "pathItems": {
          "myPathItem": {
            "post": {
              "requestBody": {
                "required": true,
                "content": {
                  "application/json": {
                    "schema": {
                      "type": "object",
                      "properties": {
                        "type": {
                          "type": "string"
                        },
                        "int": {
                          "type": "integer",
                          "exclusiveMaximum": 100,
                          "exclusiveMinimum": 0
                        },
                        "none": {
                          "type": "null"
                        },
                        "arr": {
                          "type": "array",
                          "$comment": "Array without items keyword"
                        },
                        "either": {
                          "type": [
                            "string",
                            "null"
                          ]
                        }
                      },
                      "discriminator": {
                        "propertyName": "type",
                        "x-extension": true
                      },
                      "myArbitraryKeyword": true
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  `,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "components": {
      "pathItems": {}
    }
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "summary": "My lovely API",
      "version": "1.0.0",
      "license": {
        "name": "Apache",
        "identifier": "Apache-2.0"
      }
    },
    "components": {}
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "components": {}
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "webhooks": {}
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "paths": {}
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "paths": {
      "/": {
        "get": {}
      }
    }
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "paths": {
      "/{var}": {}
    }
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "paths": {},
    "components": {
      "schemas": {
        "model": {
          "type": "object",
          "properties": {
            "one": {
              "description": "type array",
              "type": [
                "integer",
                "string"
              ]
            },
            "two": {
              "description": "type 'null'",
              "type": "null"
            },
            "three": {
              "description": "type array including 'null'",
              "type": [
                "string",
                "null"
              ]
            },
            "four": {
              "description": "array with no items",
              "type": "array"
            },
            "five": {
              "description": "singular example",
              "type": "string",
              "examples": [
                "exampleValue"
              ]
            },
            "six": {
              "description": "exclusiveMinimum true",
              "exclusiveMinimum": 10
            },
            "seven": {
              "description": "exclusiveMinimum false",
              "minimum": 10
            },
            "eight": {
              "description": "exclusiveMaximum true",
              "exclusiveMaximum": 20
            },
            "nine": {
              "description": "exclusiveMaximum false",
              "maximum": 20
            },
            "ten": {
              "description": "nullable string",
              "type": [
                "string",
                "null"
              ]
            },
            "eleven": {
              "description": "x-nullable string",
              "type": [
                "string",
                "null"
              ]
            },
            "twelve": {
              "description": "file/binary"
            }
          }
        }
      }
    }
  }`,
		`{
    "openapi": "3.1.0",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "paths": {},
    "servers": [
      {
        "url": "/v1",
        "description": "Run locally."
      },
      {
        "url": "https://production.com/v1",
        "description": "Run on production server."
      }
    ]
  }`, `{
    "openapi": "3.1.1",
    "info": {
      "title": "API",
      "version": "1.0.0"
    },
    "components": {
      "schemas": {
        "anything_boolean": true,
        "nothing_boolean": false,
        "anything_object": {},
        "nothing_object": {
          "not": {}
        }
      }
    }
  }`,
	}

	assert := require.New(t)

	for _, d := range pass {
		data := []byte(d)
		// checking json
		var o openapi.OpenAPI
		err := json.Unmarshal(data, &o)
		assert.NoError(err)
		b, err := json.MarshalIndent(o, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, b) {
			fmt.Println(string(data), "\n------------------------\n", string(b))
		}

		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))

		// testing validation
		err = o.Validate()
		assert.NoError(err)
		err = openapi.Validate(data)
		assert.NoError(err)
		// testing yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.OpenAPI
		err = yaml.Unmarshal(y, &yo)
		assert.NoError(err)
		yb, err := json.MarshalIndent(yo, "", "  ")
		assert.NoError(err)
		if !jsonpatch.Equal(data, yb) {
			fmt.Println(string(data), "\n------------------------\n", string(yb))
		}
		assert.True(jsonpatch.Equal(data, yb), cmpjson.Diff(data, yb))

	}

	fail := []string{
		`{
      "openapi": "3.1.1",
      "info": {
        "title": "API",
        "version": "1.0.0"
      },
      "components": {
        "schemas": {
          "invalid_null": null,
          "invalid_number": 0,
          "invalid_array": []
        }
      }
    }`,
		`{
      "openapi": "3.1.0",
      "info": {
        "title": "API",
        "version": "1.0.0"
      }
    }`,
		`{
      "openapi": "3.1.0",
      "info": {
        "title": "API",
        "version": "1.0.0"
      },
      "servers": [
        {
          "url": "https://example.com/{var}",
          "variables": {
            "var": {
              "enum": [],
              "default": "a"
            }
          }
        }
      ],
      "components": {}
    }`,
		`{
      "openapi": "3.1.0",
      "info": {
        "title": "API",
        "version": "1.0.0"
      },
      "paths": {},
      "servers": {
        "url": "/v1",
        "description": "Run locally."
      }
    }`,
		`{
      "openapi": "3.1.0",
      "info": {
        "title": "API",
        "version": "1.0.0"
      },
      "overlays": {}
    }`,
	}
	for _, d := range fail {
		data := []byte(d)

		err := openapi.Validate(data)
		assert.Error(err)
	}
}
