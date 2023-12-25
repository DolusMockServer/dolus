package loader

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISpecLoadType openapi3.T

type OpenAPISpecLoader struct {
	filename string
}

var _ Loader[OpenAPISpecLoadType] = &OpenAPISpecLoader{}

func NewOpenApiSpecLoader(filename string) *OpenAPISpecLoader {
	return &OpenAPISpecLoader{
		filename: filename,
	}
}

func (osl *OpenAPISpecLoader) Load() (*OpenAPISpecLoadType, error) {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(osl.filename)
	if err != nil {
		return nil, err
	}

	if err := doc.Validate(ctx); err != nil {
		return nil, err
	}

	return (*OpenAPISpecLoadType)(doc), nil
}
