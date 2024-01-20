package builder

import (
	"strconv"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/loader"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type OpenApiExpectationBuilder struct {
	loader         loader.Loader[loader.OpenAPISpecLoadType]
	fieldGenerator generator.Generator
}

var _ ExpectationBuilder = &OpenApiExpectationBuilder{}

func NewOpenApiExpectationBuilder(
	file string,
	fieldGenerator generator.Generator,
) *OpenApiExpectationBuilder {
	return &OpenApiExpectationBuilder{
		loader:         loader.NewOpenApiSpecLoader(file),
		fieldGenerator: fieldGenerator,
	}
}

func (oeb *OpenApiExpectationBuilder) BuildExpectations() (*Output, error) {
	if s, err := oeb.loader.Load(); err != nil {
		return nil, err
	} else {
		return oeb.buildExpectationsFromOpenApiSpec(s), nil
	}
}

func getRequestParameterProperty(operation *openapi3.Operation) (requestParameterProperty schema.RequestParameterProperty) {
	requestParameterProperty.PathParameterProperties = make(schema.ParameterProperties)
	requestParameterProperty.QueryParameterProperties = make(schema.ParameterProperties)

	for _, parameter := range operation.Parameters {
		parameterType := parameter.Value.In
		name := parameter.Value.Name
		required := parameter.Value.Required
		// parameterValueType := // TODO: look at paramRef.Value.Schema.Value.Type to figure this out
		if parameterType == "path" {
			requestParameterProperty.PathParameterProperties[name] = &schema.ParameterProperty{
				Required: required,
				// Type:     reflect.TypeOf(parameter.Schema),
			}
		} else if parameterType == "query" {
			requestParameterProperty.QueryParameterProperties[name] = &schema.ParameterProperty{
				Required: required,
				// Type:     reflect.TypeOf(parameter.Schema),
			}
		}
	}
	return
}

func (oeb *OpenApiExpectationBuilder) buildExpectationsFromOpenApiSpec(
	spec *loader.OpenAPISpecLoadType,
) *Output {
	var expectations []expectation.Expectation
	routeProperties := make(schema.RouteProperties)
	for path := range spec.Paths {
		refinedPath := schema.PathFromOpenApiPath(path)
		for method, operation := range spec.Paths[path].Operations() {

			routeProperties[schema.Route{
				Path:   refinedPath,
				Method: method,
			}] = getRequestParameterProperty(operation)

			for code, ref := range operation.Responses {
				if path != "/store/order/{orderId}/p" || code != "200" {
					continue
				}
				// if p != "/" || code != "200" {
				// 	continue
				// }

				logger.Log.Info(path, " ", code)

				// requestSchema := schema.RequestSchemaFromOpenApi3RequestRef(ref)
				// engine must store for each path method code then check that

				status, _ := strconv.Atoi(code)
				body := dstruct.NewGeneratedStructWithConfig(
					schema.ResponseSchemaFromOpenApi3ResponseRef(
						ref,
						"application/json",
					),
					&oeb.fieldGenerator,
				)

				expectations = append(expectations, expectation.Expectation{
					Priority: 0,
					Request: expectation.Request{
						Body:   nil,
						Path:   refinedPath,
						Method: method,
					},
					Response: expectation.Response{
						Body:          body,
						GeneratedBody: body,
						Status:        status,
					},
				})

			}
		}
	}
	return &Output{
		Expectations:    expectations,
		RouteProperties: routeProperties,
	}
}
