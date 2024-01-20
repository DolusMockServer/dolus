package engine

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

// /v1/order/2/p
// /v1/order/2/p?a=2
// /v1/order//

// if multiple matches then check

type DolusExpectationEngine struct {
	cueExpectationsFiles  []string
	expectationMatcherMap map[schema.Route][]expectation.Expectation // should be a priority queue and not just a list
	expectations          []expectation.Expectation
	cueExpectations       []expectation.Expectation
	openApiExpectations   []expectation.Expectation
	ResponseSchemas       map[schema.Route]dstruct.DynamicStructModifier
	RequestSchemas        map[schema.Route]dstruct.DynamicStructModifier
	GenerationConfig      generator.GenerationConfig
	expectationRoutes     []string
	schemaMapper          schema.Mapper
	routeProperties       schema.RouteProperties
}

var _ ExpectationEngine = &DolusExpectationEngine{}

func NewDolusExpectationEngine(
	generationConfig generator.GenerationConfig,
) (dolusExpectationEngine *DolusExpectationEngine) {
	dolusExpectationEngine = &DolusExpectationEngine{}
	dolusExpectationEngine.expectationMatcherMap = make(
		map[schema.Route][]expectation.Expectation,
	)
	dolusExpectationEngine.ResponseSchemas = make(
		map[schema.Route]dstruct.DynamicStructModifier,
	)
	dolusExpectationEngine.GenerationConfig = generationConfig

	return
}

func (e *DolusExpectationEngine) AddExpectation(
	expect expectation.Expectation,
	validateExpectationSchema bool,
	expectationType expectation.ExpectationType,
) error {
	// TODO: check if exception overrides another one i.e has same request matcher

	if validateExpectationSchema {
		if err := e.validateExpectationSchema(&expect); err != nil {
			return err
		}
	}
	// get the path with no	path parameters and query paramters
	route, err := expect.Request.RouteWithParsedPath()
	if err != nil {
		return err
	}

	e.expectationMatcherMap[*route] = append(e.expectationMatcherMap[*route], expect)

	if expectationType == expectation.Cue {
		e.cueExpectations = append(e.cueExpectations, expect)
	} else if expectationType == expectation.OpenAPI {
		e.openApiExpectations = append(e.openApiExpectations, expect)
	}

	return nil
}

func (e *DolusExpectationEngine) validateExpectationSchema(exp *expectation.Expectation) error {
	matchingResponseSchema, err := e.getMatchingResponseSchemaForRoute(exp)
	if err != nil {
		return err
	}
	expectationFieldErrors := validateExpectationResponseSchema(
		exp.Response.GeneratedBody,
		matchingResponseSchema,
	)

	if len(expectationFieldErrors) > 0 {
		return expectation.ExpectationError{
			Expectation:            *exp,
			ExpectationFieldErrors: expectationFieldErrors,
		}
	}
	return nil
}

func (e *DolusExpectationEngine) AddResponseSchemaForRoute(
	route schema.Route,
	responseSchema dstruct.DynamicStructModifier,
) error {
	if e.ResponseSchemas[route] != nil {
		return fmt.Errorf("response schema already exists for... ")
	}
	e.ResponseSchemas[route] = responseSchema
	return nil
}

func (e *DolusExpectationEngine) GetAllExpectations() map[schema.Route][]expectation.Expectation {
	return e.expectationMatcherMap
}

func (e *DolusExpectationEngine) GetExpectationRoutes() []schema.Route {
	//	return e.expectationRoutes
	//
	// TODO: implement
	return nil
}

func (e *DolusExpectationEngine) GetExpectation(
	route schema.Route,
) []expectation.Expectation {
	return e.expectationMatcherMap[route]
}

// getMatchingResponseSchemaForRoute returns the response schema for the given route
func (e *DolusExpectationEngine) getMatchingResponseSchemaForRoute(exp *expectation.Expectation) (dstruct.DynamicStructModifier, error) {
	expectationRoute := exp.Request.Route()
	parsedURL, err := url.Parse(expectationRoute.Path)
	if err != nil {
		return nil, err
	}
	for schemaRoute, responseSchema := range e.ResponseSchemas {
		if schemaRoute.Method == expectationRoute.Method {
			if pathParams, ok := schemaRoute.Match(schema.PathFromOpenApiPath(parsedURL.Path)); ok {
				// TODO: move this to getMatchingRequestSchemaForRoute
				exp.AddRequestParameterMatchers(pathParams, parsedURL.Query())
				if err := exp.ValidateRequestParameters(e.routeProperties[schemaRoute]); err != nil {
					return nil, err
				}
				return responseSchema, nil
			}
		}
	}

	return nil, fmt.Errorf(
		"no matching schema found for path=%s and operation=%s",
		expectationRoute.Path,
		expectationRoute.Method,
	)
}

func (e *DolusExpectationEngine) getMatchingRequestSchemaForRoute(
	exp *expectation.Expectation) (dstruct.DynamicStructModifier, error) {
	// TODO: implement

	// for schemaRoute, requestSchema := range e.RequestSchemas {

	// }

	return nil, fmt.Errorf("not implemented")
}

func addFieldDoesNotExistError(
	fieldName string,
	expectationFieldErrors []expectation.ExpectationFieldError,
) []expectation.ExpectationFieldError {
	fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
	return append(expectationFieldErrors, expectation.ExpectationFieldError{
		FieldName: fieldName,
		Err:       fmt.Errorf("field does not exist in the schema"),
	})
}

func addFieldMissingError(
	fieldName string,
	expectationFieldErrors []expectation.ExpectationFieldError,
) []expectation.ExpectationFieldError {
	fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
	return append(expectationFieldErrors, expectation.ExpectationFieldError{
		FieldName: fieldName,
		Err:       fmt.Errorf("required field is missing in the schema"),
	})
}

func addIncompatibleTypesError(
	fieldName string,
	expectationFieldErrors []expectation.ExpectationFieldError,
	schemaType reflect.Type,
	expectationType reflect.Type,
) []expectation.ExpectationFieldError {
	fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
	return append(expectationFieldErrors, expectation.ExpectationFieldError{
		FieldName: fieldName,
		Err: fmt.Errorf(
			"incompatible types. '%s' field is defined as type '%s' in schema but in expectation is defined as type '%s' ",
			fieldName,
			schemaType,
			expectationType,
		),
	})
}

func validateExpectationResponseSchema(
	expect dstruct.DynamicStructModifier,
	schema dstruct.DynamicStructModifier,
) (expectationFieldErrors []expectation.ExpectationFieldError) {
	// e := dstruct.New(expectation.GetSchema())
	expectationFields := expect.GetFields()
	for field, value := range schema.GetFields() {
		schemaFieldType := value.GetType()
		expectationFieldType := expectationFields[field].GetType()

		if expectationFields[field].GetFieldName() == "" && value.GetTag("required") == "true" {
			expectationFieldErrors = addFieldMissingError(field, expectationFieldErrors)
		} else if expectationFields[field].GetFieldName() != "" {
			if schemaFieldType.Kind() == reflect.Ptr {
				schemaFieldType = reflect.TypeOf(dreflect.GetUnderlyingPointerValue(value.GetValue()))
			}
			if schemaFieldType.Kind() != reflect.Struct {
				if !schemaFieldType.ConvertibleTo(expectationFieldType) || !expectationFieldType.ConvertibleTo(schemaFieldType) {
					expectationFieldErrors = addIncompatibleTypesError(field, expectationFieldErrors, schemaFieldType, expectationFieldType)
				}
			}
			if schemaFieldType.Kind() == reflect.String {
				// validate enum type
				if value.GetTag("enum") != "" {
					enumCount, _ := strconv.Atoi(value.GetTag("enum"))
					found := false
					var enumValues string
					for i := 1; i <= enumCount; i++ {
						enumValue := value.GetTag(fmt.Sprintf("enum_%d", i))
						enumValues += fmt.Sprintf("%d.'%s' ", i, enumValue)
						if !found {
							if enumValue == expectationFields[field].GetValue() {
								found = true
							}
						}
					}

					if !found {
						expectationFieldErrors = append(expectationFieldErrors, expectation.ExpectationFieldError{
							FieldName: field,
							Err: fmt.Errorf("invalid value '%s' for enum field: valid types are: \n%s",
								expectationFields[field].GetValue(), enumValues),
						})
					}
				}
			}
		}

	}

	// check for extra fields
	schemaFields := schema.GetFields()
	for field := range expect.GetFields() {
		if schemaFields[field].GetFieldName() == "" {
			// dstruct.ExtendStruct(expectation).RemoveField(field)
			expectationFieldErrors = addFieldDoesNotExistError(field, expectationFieldErrors)
		}
	}
	return
}

func (e *DolusExpectationEngine) getExpectationsForRequest(
	pathTemplate string,
	request *http.Request,
	requestParameters schema.RequestParameters,
) []expectation.Expectation {
	// check for exact matches (with query parameters)
	expectations := e.expectationMatcherMap[schema.Route{
		Path:   request.RequestURI,
		Method: request.Method,
	}]

	expectations = matchRequestParameters(requestParameters, expectations)

	if len(expectations) > 0 {
		return expectations
	}

	expectations = e.expectationMatcherMap[schema.Route{
		Path:   pathTemplate,
		Method: request.Method,
	}]

	expectations = matchRequestParameters(requestParameters, expectations)

	return expectations

}

func matchRequestParameters(
	requestParameters schema.RequestParameters,
	expectations []expectation.Expectation,
) (filtered []expectation.Expectation) {
	// TODO: get path paratemers of request
	// if expectation already has a pathe parameter ignore the additional filter ones
	// requestPathParameters = request.URL.Query()

	for _, e := range expectations {
		if e.Request.Parameters == nil || e.Request.Parameters.Match(requestParameters) {
			filtered = append(filtered, e)
		}
	}
	return
}

// GetResponseForRequest returns the response for the given request
func (e *DolusExpectationEngine) GetResponseForRequest(
	request *http.Request,
	requestParameters schema.RequestParameters,
	pathTemplate string,
) (*expectation.Response, error) {
	expectations := e.getExpectationsForRequest(pathTemplate, request, requestParameters)

	// TODO: find the right expectation depending on request matchers and priority
	// now sift through the expectations and look at path paraem
	if len(expectations) == 0 {
		return nil, fmt.Errorf("no expectation found for path and HTTP method")
	}

	// findMatchingPathParameters(expectations)
	// findHighestPriorityExpectation(expectations)
	currentExpectation := expectations[0]
	for _, v := range expectations {
		if request.URL.Path == v.Request.Path {
			currentExpectation = v
			return &currentExpectation.Response, nil
		}
		if v.Priority > currentExpectation.Priority {
			currentExpectation = v
		}
	}

	return &currentExpectation.Response, nil
}

// GetCueExpectations returns the cue expectations
func (e *DolusExpectationEngine) GetCueExpectations() expectation.Expectations {
	// TODO: implement
	return expectation.Expectations{
		Expectations: e.cueExpectations,
	}
}

func (e *DolusExpectationEngine) SetRouteProperties(routeProperties schema.RouteProperties) {
	e.routeProperties = routeProperties
}
