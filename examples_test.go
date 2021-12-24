package openapi_test

import (
	"embed"
	"encoding/json"
	"log"

	"github.com/chanced/openapi"
)

func ExampleUsingLoaderToOpenExistingAPI() {
	//go:embed testdata/examples/openapi.yaml
	var files embed.FS
	d, err := files.Open("openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}
	o, err := openapi.Load(d, openapi.NewResolver(
		openapi.Openers{
			"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.1": &openapi.HTTPOpener{},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// o is now a *ResolvedOpenAPI
	_ = o
}

func ExampleUnmarshalingOpenAPI() {
	j := `{
		"openapi": "3.1.0",
		"info": {
		  "title": "Webhook Example",
		  "version": "1.0.0"
		},
		"webhooks": {
		  "newPet": {
			"post": {
			  "requestBody": {
				"description": "Information about a new pet in the system",
				"content": {
				  "application/json": {
					"schema": {
					  "$ref": "#/components/schemas/Pet"
					}
				  }
				}
			  },
			  "responses": {
				"200": {
				  "description": "Return a 200 status to indicate that the data was received successfully"
				}
			  }
			}
		  }
		},
		"components": {
		  "schemas": {
			"Pet": {
			  "required": [
				"id",
				"name"
			  ],
			  "properties": {
				"id": {
				  "type": "integer",
				  "format": "int64"
				},
				"name": {
				  "type": "string"
				},
				"tag": {
				  "type": "string"
				}
			  }
			}
		  }
		}
	  }`
	var o *openapi.OpenAPI
	if err := json.Unmarshal([]byte(j), &o); err != nil {
		log.Fatal(err)
	}
}
