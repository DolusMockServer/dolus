package core

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ResponseSchema struct {
	Schema
	StatusCode string
}

var _ schema = &ResponseSchema{}

func NewResponseSchemaFromOpenApi3Ref(path, method, statusCode string, ref *openapi3.ResponseRef, mediaType string,
) *ResponseSchema {

	return &ResponseSchema{
		Schema: Schema{
			Path:   path,
			Method: method,
			schema: getSchemaFromOpenApi3Spec(ref, mediaType),
		},
		StatusCode: statusCode,
	}

}

func NewResponseSchemaFromAny(path, method, statusCode string, config any) *ResponseSchema {

	return &ResponseSchema{
		Schema: Schema{
			Path:   path,
			Method: method,
			schema: getStructFromAny(config),
		},
		StatusCode: statusCode,
	}
}
