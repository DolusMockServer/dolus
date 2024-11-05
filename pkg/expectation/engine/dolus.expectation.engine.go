package engine

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/matcher"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type DolusExpectationEngine struct {
	expectationMatcherMap map[schema.Route][]expectation.Expectation // should be a priority queue and not just a list
	expectations          []expectation.Expectation
	GenerationConfig      generator.GenerationConfig
	routeManager          RouteManager
}

var _ ExpectationEngine = &DolusExpectationEngine{}

func NewDolusExpectationEngine(
	generationConfig generator.GenerationConfig,
	routeManager RouteManager,
	expectations []expectation.Expectation,
) (dolusExpectationEngine *DolusExpectationEngine) {
	dolusExpectationEngine = &DolusExpectationEngine{}
	dolusExpectationEngine.expectationMatcherMap = make(
		map[schema.Route][]expectation.Expectation,
	)
	dolusExpectationEngine.GenerationConfig = generationConfig
	dolusExpectationEngine.routeManager = routeManager
	dolusExpectationEngine.AddExpectations(expectations)

	return
}

func (e *DolusExpectationEngine) AddExpectations(expectations []expectation.Expectation) {
	for _, exp := range expectations {
		if err := e.AddExpectation(exp); err != nil {
			slog.Debug(fmt.Sprintf("Error adding expectation:\n%s\n", err.Error()), "expectation", exp)
		}
	}
}

func (e *DolusExpectationEngine) AddExpectation(
	expect expectation.Expectation) error {
	// TODO: check if exception overrides another one i.e has same request matcher

	// use the request path to add query parameters to the expectation and override existing query parameters
	addQueryParameterMatcher(&expect)

	if expect.ExpectationType == expectation.Custom {
		if err := e.validateExpectationSchema(&expect); err != nil {
			return err
		}
	}
	// // get the path with no	path parameters and query parameters
	route, err := expect.Request.RouteWithParsedPath()
	if err != nil {
		return err
	}

	e.expectationMatcherMap[*route] = append(e.expectationMatcherMap[*route], expect)

	e.expectations = append(e.expectations, expect)

	return nil
}

// addQueryParameterMatcher adds query parameter matchers from the expectation's request path to the expectation. If the expectation already has query parameters, they will be overwritten.
func addQueryParameterMatcher(e *expectation.Expectation) error {
	parsedURL, err := url.Parse(e.Request.Path)
	if err != nil {
		return fmt.Errorf(
			"failed to add query parameters for expectation with path '%s': %w",
			e.Request.Path,
			err,
		)
	}
	queryParams := parsedURL.Query()
	if e.Request.Parameters == nil {
		e.Request.Parameters = &expectation.RequestParameters{}
	}
	if e.Request.Parameters.Query == nil {
		e.Request.Parameters.Query = make(map[string]any)
	}
	for k, v := range queryParams {
		value := v
		e.Request.Parameters.Query[k] = matcher.NewStringArrayMatcher(&value, "eq")
	}
	return nil
}

func (e *DolusExpectationEngine) validateExpectationSchema(exp *expectation.Expectation) error {
	matchingResponseSchema, err := e.getResponseSchemaForExpectation(exp)
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

// GetResponseForRequest returns the response for the given request
func (e *DolusExpectationEngine) GetResponseForRequest(
	request *http.Request,
	requestParameters schema.RequestParameters,
	pathTemplate string,
) (*expectation.Response, error) {
	// TODO: if Route does not exist in route manager return an error
	expectations := e.getExpectationsForRequest(pathTemplate, request, requestParameters)

	// TODO: find the right expectation depending on request matchers and priority
	// now sift through the expectations and look at path paraem
	if len(expectations) == 0 {
		return nil, fmt.Errorf("no expectation found for path and HTTP method")
	}

	// findHighestPriorityExpectation(expectations)
	currentExpectation := expectations[0]
	for _, v := range expectations {
		if v.Priority > currentExpectation.Priority {
			currentExpectation = v
		}
	}

	return &currentExpectation.Response, nil
}

func (e *DolusExpectationEngine) GetExpectations(
	expectationType *expectation.ExpectationType,
	path *string,
	method *string,
) []expectation.Expectation {
	var filteredExpectations []expectation.Expectation

	for _, exp := range e.expectations {
		if path != nil && *path != exp.Request.Path {
			continue
		}
		if method != nil && *method != exp.Request.Method {
			continue
		}
		if expectationType != nil && *expectationType != exp.ExpectationType {
			continue
		}
		filteredExpectations = append(filteredExpectations, exp)
	}

	return filteredExpectations
}

// GetRoutes implements ExpectationEngine.
func (e *DolusExpectationEngine) GetRoutes() []schema.Route {
	return e.routeManager.GetRoutes()
}

// getMatchingResponseSchemaForRoute returns the response schema for the given route
// get the response schema for the expectation if it has one
func (e *DolusExpectationEngine) getResponseSchemaForExpectation(
	exp *expectation.Expectation,
) (dstruct.DynamicStructModifier, error) {
	expectationRoute := exp.Request.Route()
	parsedURL, err := url.Parse(expectationRoute.Path)
	if err != nil {
		return nil, err
	}
	for schemaRoute, routeProperty := range e.routeManager.GetRouteProperties() {
		if schemaRoute.Method == expectationRoute.Method {
			if pathParams, ok := schemaRoute.Match(schema.PathFromOpenApiPath(parsedURL.Path)); ok {
				// TODO: move this to getMatchingRequestSchemaForRoute

				// we can add the path parameters here as we can only figure out what the path parameters are
				// once know what the route is
				if err := addPathParametersFromRequestPath(exp, pathParams); err != nil {
					return nil, fmt.Errorf(
						"failed to add path parameters for expectation with path %s: %w",
						exp.Request.Path,
						err,
					)
				}
				// found the matching right path and operation now validate the request parameter
				// TODO: only validate query parameters if not generic path
				if err := validateRequestParameters(exp, routeProperty.RequestParameterProperty); err != nil {
					return nil, err
				}
				return routeProperty.ResponseSchema, nil
			}
		}
	}

	return nil, fmt.Errorf(
		"no matching schema found for path=%s and operation=%s",
		expectationRoute.Path,
		expectationRoute.Method,
	)
}

// TODO  refractor add and validate code
func addPathParametersFromRequestPath(
	e *expectation.Expectation,
	pathParams map[string]string,
) error {
	if e.Request.Parameters == nil {
		e.Request.Parameters = &expectation.RequestParameters{}
	}

	return addPathParameters(pathParams, e)
}

func addPathParameters(pathParams map[string]string, e *expectation.Expectation) error {
	if e.Request.Parameters.Path == nil {
		e.Request.Parameters.Path = make(map[string]any)
	}
	for k, v := range pathParams {
		matchType := "eq"
		value := v
		if v == ":"+k {
			matchType = "has"
		} else if strings.TrimSpace(v) == "" {
			return fmt.Errorf("path parameter '%s' is empty", k)
		}
		e.Request.Parameters.Path[k] = matcher.NewStringMatcher(&value, matchType)

	}
	return nil
}

func validateRequestParameters(
	expectation *expectation.Expectation,
	routeProperty schema.RequestParameterProperty,
) error {
	// Validate Path and Query Parameters
	if err := checkParametersExistence("Path", routeProperty.PathParameterProperties, expectation.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkRequiredPathParameters(routeProperty.PathParameterProperties, expectation.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkParametersExistence("Query", routeProperty.QueryParameterProperties, expectation.Request.Parameters.Query); err != nil {
		return err
	}

	if err := checkRequiredQueryParameters(routeProperty.QueryParameterProperties, expectation.Request.Parameters.Query); err != nil {
		return err
	}
	return nil
}

// checkParametersExistence checks for extra parameters not defined in the schema
func checkParametersExistence(
	paramType string,
	properties schema.ParameterProperties,
	parameters map[string]any,
) error {
	for name := range parameters {
		if properties[name] == nil {
			return fmt.Errorf("%s parameter '%s' does not exist", paramType, name)
		}
	}
	return nil
}

func checkRequiredQueryParameters(
	properties schema.ParameterProperties,
	parameters map[string]any,
) error {
	for name, param := range properties {
		if param.Required &&
			(parameters[name] == nil || len(*parameters[name].(*matcher.StringArrayMatcher).GetValue()) == 0) {
			return fmt.Errorf("required query parameter '%s' is missing", name)
		}
	}
	return nil
}

func checkRequiredPathParameters(
	properties schema.ParameterProperties,
	parameters map[string]any,
) error {
	for name, param := range properties {
		if param.Required &&
			(parameters[name] == nil || *parameters[name].(*matcher.StringMatcher).GetValue() == "") {
			return fmt.Errorf("required path parameter '%s' is missing", name)
		}
	}
	return nil
}

func (e *DolusExpectationEngine) getMatchingRequestSchemaForRoute(
	exp *expectation.Expectation,
) (dstruct.DynamicStructModifier, error) {
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
	var expectationFields dstruct.FieldData
	if expect == nil {
		expectationFields = make(dstruct.FieldData)
	} else {
		expectationFields = expect.GetFields()
	}
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
	for field := range expectationFields {
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
	if expectations := e.findExpectationMatches(request.URL.Path, request, requestParameters); len(
		expectations,
	) > 0 {
		return expectations
	}

	return e.findExpectationMatches(pathTemplate, request, requestParameters)
}

// returns expectations for a specific route
func (e *DolusExpectationEngine) findExpectationMatches(
	requestPath string,
	request *http.Request,
	requestParameters schema.RequestParameters,
) (filtered []expectation.Expectation) {
	expectations := e.expectationMatcherMap[schema.Route{
		Path:   requestPath,
		Method: request.Method,
	}]

	for _, expectation := range expectations {
		requestMatcher := matcher.NewRequestMatcher(expectation.Request)
		if requestMatcher.Matches(request, &requestParameters) {
			filtered = append(filtered, expectation)
		}
	}
	return
}
