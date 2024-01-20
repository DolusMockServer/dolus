package schema

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type RequestSchema struct {
	body any
}

func ResponseSchemaFromOpenApi3ResponseRef(
	ref *openapi3.ResponseRef,
	mediaType string,
) any {
	return getSchemaFromOpenApi3Spec(ref, mediaType)
}

func RequestSchemaFromOpenApi3RequestRef(
	ref *openapi3.RequestBody,
	mediaType string,
) any {
	// TODO: implement
	return nil
}

func SchemaFromAny(config any) any {
	return getStructFromAny(config)
}

// TODO: look into Callbacks from the openapi3 spec

func PathFromOpenApiPath(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}
