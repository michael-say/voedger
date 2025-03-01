/* PROMPT:

The application components are described with IAppDef interface.


I want to have an implementation of a function `GenOpenAPISchema`, which generates OpenAPI v3 JSON Schema of the API to my application:
```go
type PublishedTypeCallback func(appdef.IType, []OperationKind) error

type IAppDefWithPublishedTypes interface {
	IAppDef
        // PublishedTypes lists resources available to the published role in the workspace and ancestors (including resources available to non-authenticated requests):
        // - Documents (appdef.IODoc, appdef.IORecord, appdef.ICDoc, appdef.ICRecord, appdef.IWDoc, appdef.IWRecord)
        // - Views (appdef.IView)
        // - Commands (appdef.ICommand)
        // - Queries (appdef.IQuery)
    PublishedTypes(ws appdef.IWorkspace, role appdef.QName, callback PublishedTypeCallback) error
}

func GenOpenAPISchema(appDef IAppDefWithPublishedTypes, ws appdef.IWorkspace, role appdef.QName, writer io.Wrter, schemaTitle, schemaVerstion string) error
```

Paths in schema are generated per published type per operation. Example path for opeation `OperationKind_Insert`:
```
# Create document or record
## Motivation
Create a new CDoc, WDoc, CRecord or WRecord using API

## Functional Design
POST `/api/v2/users/{owner}/apps/{app}/workspaces/{wsid}/docs/{pkg}.{table}`

### Headers
| Key | Value |
| --- | --- |
| Authorization | Bearer {PrincipalToken} |
| Content-type | application/json |

### Parameters
owner (string) - name of a user who owns the application
app (string) - name of an application
wsid (int64) - the ID of workspace
pkg, table (string) - identifies a table (document or record)

### Request
JSON of the CDoc/WDoc/CRecord/WRecord


### Result
| Code | Description | Body
| --- | --- | --- |
| 201 | Created | current WLog offset and the new IDs, see below |
| 400 | Bad Request, e.g. Record requires sys.ParentID | error object, see below |
| 401 | Unauthorized | error object, see below |
| 403 | Forbidden | error object, see below|
| 404 | Table Not Found | error object, see below |
| 405 | Method Not Allowed, table is an ODoc/ORecord | error object, see below |

Example result 201:
```json
{
    "CurrentWLogOffset":114,
    "NewIDs": {
        "1":322685000131212
    }
}

Example non 2xx result (error object):
{
  "code": 105,
  "error": "invalid field name: bl!ng"
}
```

Create an implementation of a function in the file `openapi2.go`, which describes path for `OperationKind_Insert` for documents.

*/

package appdef

import (
	"encoding/json"
	"fmt"
	"io"
)

type PublishedTypeCallback func(IType, []OperationKind) error

// OpenAPI schema structures
type openAPISchema struct {
	OpenAPI string               `json:"openapi"`
	Info    openAPIInfo          `json:"info"`
	Paths   map[string]*pathItem `json:"paths"`
}

type openAPIInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type pathItem struct {
	Post *operation `json:"post,omitempty"`
}

type operation struct {
	Tags        []string             `json:"tags"`
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	Parameters  []parameter          `json:"parameters"`
	RequestBody *requestBody         `json:"requestBody"`
	Responses   map[string]*response `json:"responses"`
}

type parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Schema      *schema `json:"schema"`
}

type requestBody struct {
	Required bool                  `json:"required"`
	Content  map[string]*mediaType `json:"content"`
}

type response struct {
	Description string                `json:"description"`
	Content     map[string]*mediaType `json:"content"`
}

type mediaType struct {
	Schema *schema `json:"schema"`
}

type schema struct {
	Type                 string             `json:"type,omitempty"`
	Format               string             `json:"format,omitempty"`
	Properties           map[string]*schema `json:"properties,omitempty"`
	AdditionalProperties *schema            `json:"additionalProperties,omitempty"`
}

type IAppDefWithPublishedTypes interface {
	IAppDef
	PublishedTypes(ws IWorkspace, role QName, callback PublishedTypeCallback) error
}

// OpenAPI schema generation
func GenOpenAPISchema(appDef IAppDefWithPublishedTypes, ws IWorkspace, role QName, writer io.Writer, schemaTitle, schemaVersion string) error {
	schema := &openAPISchema{
		OpenAPI: "3.0.0",
		Info: openAPIInfo{
			Title:   schemaTitle,
			Version: schemaVersion,
		},
		Paths: make(map[string]*pathItem),
	}

	err := appDef.PublishedTypes(ws, role, func(t IType, ops []OperationKind) error {
		for _, op := range ops {
			switch op {
			case OperationKind_Insert:
				if err := addInsertPath(t, schema); err != nil {
					return fmt.Errorf("error adding insert path for %s: %w", t.QName(), err)
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error processing published types: %w", err)
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(schema)
}

func addInsertPath(t IType, aschema *openAPISchema) error {
	// Only handle document/record types
	switch t.Kind() {
	case TypeKind_CDoc, TypeKind_WDoc, TypeKind_CRecord, TypeKind_WRecord:
		// Continue processing
	default:
		return nil
	}

	pathTemplate := "/api/v2/users/{owner}/apps/{app}/workspaces/{wsid}/docs/%s.%s"
	path := fmt.Sprintf(pathTemplate, t.QName().Pkg(), t.QName().Entity())

	pathItem := &pathItem{
		Post: &operation{
			Tags:        []string{"Documents"},
			Summary:     fmt.Sprintf("Create new %s", t.QName()),
			Description: fmt.Sprintf("Create a new instance of %s", t.QName()),
			Parameters: []parameter{
				{
					Name:        "owner",
					In:          "path",
					Required:    true,
					Schema:      &schema{Type: "string"},
					Description: "Name of user who owns the application",
				},
				{
					Name:        "app",
					In:          "path",
					Required:    true,
					Schema:      &schema{Type: "string"},
					Description: "Name of application",
				},
				{
					Name:        "wsid",
					In:          "path",
					Required:    true,
					Schema:      &schema{Type: "integer", Format: "int64"},
					Description: "Workspace ID",
				},
			},
			RequestBody: &requestBody{
				Required: true,
				Content: map[string]*mediaType{
					"application/json": {
						Schema: generateRequestSchema(t),
					},
				},
			},
			Responses: map[string]*response{
				"201": {
					Description: "Created successfully",
					Content: map[string]*mediaType{
						"application/json": {
							Schema: &schema{
								Type: "object",
								Properties: map[string]*schema{
									"CurrentWLogOffset": {
										Type: "integer",
									},
									"NewIDs": {
										Type: "object",
										AdditionalProperties: &schema{
											Type: "integer",
										},
									},
								},
							},
						},
					},
				},
				"400": generateErrorResponse("Bad Request"),
				"401": generateErrorResponse("Unauthorized"),
				"403": generateErrorResponse("Forbidden"),
				"404": generateErrorResponse("Not Found"),
				"405": generateErrorResponse("Method Not Allowed"),
			},
		},
	}

	aschema.Paths[path] = pathItem
	return nil
}

func generateRequestSchema(t IType) *schema {
	props := make(map[string]*schema)

	// Add fields if type implements IWithFields
	if fieldsType, ok := t.(IWithFields); ok {
		for _, field := range fieldsType.Fields() {
			fieldSchema := &schema{}

			switch field.DataKind() {
			case DataKind_int32:
				fieldSchema.Type = "integer"
				fieldSchema.Format = "int32"
			case DataKind_int64:
				fieldSchema.Type = "integer"
				fieldSchema.Format = "int64"
			case DataKind_float32:
				fieldSchema.Type = "number"
				fieldSchema.Format = "float"
			case DataKind_float64:
				fieldSchema.Type = "number"
				fieldSchema.Format = "double"
			case DataKind_string:
				fieldSchema.Type = "string"
			case DataKind_bool:
				fieldSchema.Type = "boolean"
			}

			props[field.Name()] = fieldSchema
		}
	}

	return &schema{
		Type:       "object",
		Properties: props,
	}
}

func generateErrorResponse(description string) *response {
	return &response{
		Description: description,
		Content: map[string]*mediaType{
			"application/json": {
				Schema: &schema{
					Type: "object",
					Properties: map[string]*schema{
						"code": {
							Type: "integer",
						},
						"error": {
							Type: "string",
						},
					},
				},
			},
		},
	}
}
