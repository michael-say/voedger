/*
 * Copyright (c) 2025-present unTill Software Development Group B.V.
 * @author Michael Saigachenko
 */
package query2

import (
	_ "embed"
)

//go:embed swagger-ui.html
var swaggerUI_HTML string

const (
	errorSchemaName          = "Error"
	errorSchemaRef           = "#/components/schemas/" + errorSchemaName
	principalTokenSchemaName = "PrincipalToken"
	principalTokenSchemaRef  = "#/components/schemas/" + principalTokenSchemaName
	bearerAuth               = "BearerAuth"
	authenticationTag        = "Authentication"
)

const (
	methodGet    = "get"
	methodPost   = "post"
	methodPut    = "put"
	methodDelete = "delete"
	methodPatch  = "patch"
)

// Content types
const (
	applicationJSON = "application/json"
)

// Descriptions
const (
	descrOK = "OK"
)

// Status codes
const (
	statusCode200 = "200"
	statusCode400 = "400"
	statusCode401 = "401"
	statusCode403 = "403"
	statusCode429 = "429"
	statusCode500 = "500"
)

// OpenAPI schema constants
const (
	schemaTypeObject  = "object"
	schemaTypeString  = "string"
	schemaTypeInteger = "integer"
	schemaTypeNumber  = "number"
	schemaTypeBoolean = "boolean"
	schemaTypeArray   = "array"

	schemaMethodPost = "post"
	schemaMethodGet  = "get"

	schemaFormatInt32  = "int32"
	schemaFormatInt64  = "int64"
	schemaFormatFloat  = "float"
	schemaFormatDouble = "double"
	schemaFormatByte   = "byte"

	schemaKeyType        = "type"
	schemaKeyFormat      = "format"
	schemaKeyDescription = "description"
	schemaKeyProperties  = "properties"
	schemaKeyRequired    = "required"
	schemaKeyContent     = "content"
	schemaKeySchema      = "schema"
	schemaKeyItems       = "items"
	schemaKeyOneOf       = "oneOf"
	schemaKeyRef         = "$ref"
	schemaKeyRequestBody = "requestBody"
	schemaKeyResponses   = "responses"
	schemaKeyParameters  = "parameters"
	schemaKeySecurity    = "security"
	schemaKeyTags        = "tags"
)
