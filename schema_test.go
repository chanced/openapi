package openapi_test

import (
	"encoding/json"
	"fmt"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/require"
	yaml "sigs.k8s.io/yaml"

	"github.com/chanced/cmpjson"
	"github.com/chanced/openapi"
)

func TestSchema(t *testing.T) {
	assert := require.New(t)

	j := []string{`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"$id": "https://example.com/tree",
		"$dynamicAnchor": "node",
		"type": "object",
		"properties": {
		  "data": true,
		  "children": {
			"type": "array",
			"items": { "$dynamicRef": "#node" }
		  }
		}, 
		"discriminator": {
			"propertyName": "type",
			"x-extension": true
		}
	  }`,
		`{
		"$id": "https://example.com/person.schema.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"title": "Person",
		"type": "object",
		"properties": {
		  "firstName": {
			"type": "string",
			"description": "The person's first name."
		  },
		  "lastName": {
			"type": "string",
			"description": "The person's last name."
		  },
		  "age": {
			"description": "Age in years which must be equal to or greater than zero.",
			"type": "integer",
			"minimum": 0
		  }
		}
	  }`,
		`{
		"$ref": "#/$defs/enabledToggle",
		"default": true
	}`,
		`{
		"title": "Feature A",
		"properties": {
			"enabled": {
				"$ref": "#/$defs/enabledToggle",
				"default": true
			}
		}
	}`,
		`{
		"title": "Feature B",
		"properties": {
			"enabled": {
				"description": "If set to null, Feature B inherits the enabled value from Feature A",
				"$ref": "#/$defs/enabledToggle"
			}
		}
	}`,
		`{
		"title": "Feature list",
		"type": "array",
		"prefixObjs": [
			{
				"title": "Feature A",
				"properties": {
					"enabled": {
						"$ref": "#/$defs/enabledToggle",
						"default": true
					}
				}
			},
			{
				"title": "Feature B",
				"properties": {
					"enabled": {
						"description": "If set to null, Feature B inherits the enabled value from Feature A",
						"$ref": "#/$defs/enabledToggle"
					}
				}
			}
		],
		"$defs": {
			"enabledToggle": {
				"title": "Enabled",
				"description": "Whether the feature is enabled (true), disabled (false), or under automatic control (null)",
				"type": ["boolean", "null"],
				"default": null
			}
		}
	}`,
		`{
		"$id": "https://example.com/geographical-location.schema.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"title": "Longitude and Latitude Values",
		"description": "A geographical coordinate.",
		"required": [ "latitude", "longitude" ],
		"type": "object",
		"properties": {
		  "latitude": {
			"type": "number",
			"minimum": -90,
			"maximum": 90
		  },
		  "longitude": {
			"type": "number",
			"minimum": -180,
			"maximum": 180
		  }
		}
	  }`,
		`{
		"$id": "https://example.com/card.schema.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"description": "A representation of a person, company, organization, or place",
		"type": "object",
		"required": [ "familyName", "givenName" ],
		"properties": {
		  "fn": {
			"description": "Formatted Name",
			"type": "string"
		  },
		  "familyName": {
			"type": "string"
		  },
		  "givenName": {
			"type": "string"
		  },
		  "additionalName": {
			"type": "array",
			"items": {
			  "type": "string"
			}
		  },
		  "honorificPrefix": {
			"type": "array",
			"items": {
			  "type": "string"
			}
		  },
		  "honorificSuffix": {
			"type": "array",
			"items": {
			  "type": "string"
			}
		  },
		  "nickname": {
			"type": "string"
		  },
		  "url": {
			"type": "string"
		  },
		  "email": {
			"type": "object",
			"properties": {
			  "type": {
				"type": "string"
			  },
			  "value": {
				"type": "string"
			  }
			}
		  },
		  "tel": {
			"type": "object",
			"properties": {
			  "type": {
				"type": "string"
			  },
			  "value": {
				"type": "string"
			  }
			}
		  },
		  "adr": { "$ref": "https://example.com/address.schema.json" },
		  "geo": { "$ref": "https://example.com/geographical-location.schema.json" },
		  "tz": {
			"type": "string"
		  },
		  "photo": {
			"type": "string"
		  },
		  "logo": {
			"type": "string"
		  },
		  "sound": {
			"type": "string"
		  },
		  "bday": {
			"type": "string"
		  },
		  "title": {
			"type": "string"
		  },
		  "role": {
			"type": "string"
		  },
		  "org": {
			"type": "object",
			"properties": {
			  "organizationName": {
				"type": "string"
			  },
			  "organizationUnit": {
				"type": "string"
			  }
			}
		  }
		}
	  }`,
		`{
		"$id": "https://example.com/calendar.schema.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"description": "A representation of an event",
		"type": "object",
		"required": [ "dtstart", "summary" ],
		"properties": {
		  "dtstart": {
			"type": "string",
			"description": "Event starting time"
		  },
		  "dtend": {
			"type": "string",
			"description": "Event ending time"
		  },
		  "summary": {
			"type": "string"
		  },
		  "location": {
			"type": "string"
		  },
		  "url": {
			"type": "string"
		  },
		  "duration": {
			"type": "string",
			"description": "Event duration"
		  },
		  "rdate": {
			"type": "string",
			"description": "Recurrence date"
		  },
		  "rrule": {
			"type": "string",
			"description": "Recurrence rule"
		  },
		  "category": {
			"type": "string"
		  },
		  "description": {
			"type": "string"
		  },
		  "geo": {
			"$ref": "https://example.com/geographical-location.schema.json"
		  }
		}
	  }`,
		`{
		"$id": "https://example.com/address.schema.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"description": "An address similar to http://microformats.org/wiki/h-card",
		"type": "object",
		"properties": {
		  "post-office-box": {
			"type": "string"
		  },
		  "extended-address": {
			"type": "string"
		  },
		  "street-address": {
			"type": "string"
		  },
		  "locality": {
			"type": "string"
		  },
		  "region": {
			"type": "string"
		  },
		  "postal-code": {
			"type": "string"
		  },
		  "country-name": {
			"type": "string"
		  }
		},
		"required": [ "locality", "region", "country-name" ],
		"dependentRequired": {
		  "post-office-box": [ "street-address" ],
		  "extended-address": [ "street-address" ]
		}
	  }`,
	}
	for _, d := range j {
		var data = []byte(d)
		var v *openapi.SchemaObj
		err := json.Unmarshal(data, &v)
		assert.NoError(err)

		err = json.Unmarshal(data, &v)
		assert.NoError(err)

		b, err := json.MarshalIndent(v, "", "  ")
		assert.NoError(err)

		if !jsonpatch.Equal(data, b) {
			fmt.Println(d)
			fmt.Println(string(b))
			// litter.Dump(v)
		}
		assert.True(jsonpatch.Equal(data, b), cmpjson.Diff(data, b))
		assert.NoError(err)
		b, err = yaml.Marshal(v)
		assert.NoError(err)

		var s *openapi.SchemaObj
		err = yaml.Unmarshal(b, &s)
		assert.NoError(err)
		b, err = json.MarshalIndent(s, "", "  ")
		assert.NoError(err)
		assert.True(jsonpatch.Equal(b, data))

		// checking yaml

		y, err := yaml.JSONToYAML(data)
		assert.NoError(err)
		var yo openapi.SchemaObj
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
