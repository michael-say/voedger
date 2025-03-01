package appdef

/*
type PublishedTypeCallback func(IType, []OperationKind) error

type IAppDefWithPublishedTypes interface {
	IAppDef
	PublishedTypes(ws IWorkspace, role QName, callback PublishedTypeCallback) error
}

func GenOpenAPISchema(appDef IAppDefWithPublishedTypes, ws IWorkspace, role QName, writer io.Writer, schemaTitle, schemaVersion string) error {
	paths := make(map[string]interface{})

	err := appDef.PublishedTypes(ws, role, func(t IType, ops []OperationKind) error {
		for _, op := range ops {
			if op == OperationKind_Insert {
				path := generateInsertPath(t)
				paths[path] = generateInsertPathSchema(t)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	schema := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":   schemaTitle,
			"version": schemaVersion,
		},
		"paths": paths,
	}

	return json.NewEncoder(writer).Encode(schema)
}

func generateInsertPath(t IType) string {
	return fmt.Sprintf("/api/v2/users/{owner}/apps/{app}/workspaces/{wsid}/docs/%s", t.QName())
}

func generateInsertPathSchema(t IType) map[string]interface{} {
	return map[string]interface{}{
		"post": map[string]interface{}{
			"summary":     "Create document or record",
			"description": "Create a new CDoc, WDoc, CRecord or WRecord using API",
			"parameters": []map[string]interface{}{
				{
					"name":     "owner",
					"in":       "path",
					"required": true,
					"schema": map[string]string{
						"type": "string",
					},
				},
				{
					"name":     "app",
					"in":       "path",
					"required": true,
					"schema": map[string]string{
						"type": "string",
					},
				},
				{
					"name":     "wsid",
					"in":       "path",
					"required": true,
					"schema": map[string]string{
						"type": "integer",
					},
				},
				{
					"name":     "pkg",
					"in":       "path",
					"required": true,
					"schema": map[string]string{
						"type": "string",
					},
				},
				{
					"name":     "table",
					"in":       "path",
					"required": true,
					"schema": map[string]string{
						"type": "string",
					},
				},
			},
			"requestBody": map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]string{
							"$ref": fmt.Sprintf("#/components/schemas/%s", t.QName()),
						},
					},
				},
			},
			"responses": map[string]interface{}{
				"201": map[string]interface{}{
					"description": "Created",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"CurrentWLogOffset": map[string]string{
										"type": "integer",
									},
									"NewIDs": map[string]interface{}{
										"type": "object",
										"additionalProperties": map[string]string{
											"type": "integer",
										},
									},
								},
							},
						},
					},
				},
				"400": map[string]interface{}{
					"description": "Bad Request",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"code": map[string]string{
										"type": "integer",
									},
									"error": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"401": map[string]interface{}{
					"description": "Unauthorized",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"code": map[string]string{
										"type": "integer",
									},
									"error": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"403": map[string]interface{}{
					"description": "Forbidden",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"code": map[string]string{
										"type": "integer",
									},
									"error": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"404": map[string]interface{}{
					"description": "Table Not Found",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"code": map[string]string{
										"type": "integer",
									},
									"error": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"405": map[string]interface{}{
					"description": "Method Not Allowed",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"code": map[string]string{
										"type": "integer",
									},
									"error": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
*/
