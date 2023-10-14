package expectation

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPISpecLoadType openapi3.T

type OpenAPISpecLoader struct {
	filename string
}

func NewOpenOPISpecLoader(filename string) *OpenAPISpecLoader {
	return &OpenAPISpecLoader{
		filename: filename,
	}
}

var _ Loader[OpenAPISpecLoadType] = &OpenAPISpecLoader{}

func (osl *OpenAPISpecLoader) load() (*OpenAPISpecLoadType, error) {
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
