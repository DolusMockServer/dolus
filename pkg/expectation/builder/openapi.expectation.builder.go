package builder

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
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
	routeManager := engine.NewRouteManager()
	for path := range spec.Paths.Map() {
		refinedPath := schema.PathFromOpenApiPath(path)
		for method, operation := range spec.Paths.Map()[path].Operations() {

			for code, ref := range operation.Responses.Map() {
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

				routeManager.AddRoute(schema.Route{
					Path:   refinedPath,
					Method: method,
				}, schema.RouteProperty{
					RequestParameterProperty: getRequestParameterProperty(operation),
					RequestSchema:            nil,
					ResponseSchema:           body,
				})

				expectations = append(expectations, expectation.Expectation{
					Priority:        0,
					ExpectationType: expectation.Default,
					Request: expectation.Request{
						Body:   nil,
						Path:   refinedPath,
						Method: method,
					},
					Response: expectation.Response{
						Body:          structToMap(body),
						GeneratedBody: body,
						Status:        status,
					},
				})

			}
		}
	}
	return &Output{
		Expectations: expectations,
		RouteManager: routeManager,
	}
}

func structToMap(obj *dstruct.GeneratedStructImpl) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range obj.GetFields() {
		if strings.Contains(k, ".") {
			continue
		}
		fieldValueKind := v.GetType().Kind()
		var fieldValue interface{}

		if fieldValueKind == reflect.Struct {
			fieldValue = structToMap(dstruct.NewGeneratedStruct(v.GetValue()))
		} else {
			fieldValue = v.GetValue()

		}

		result[v.GetJsonName()] = fieldValue

	}

	return result
}
