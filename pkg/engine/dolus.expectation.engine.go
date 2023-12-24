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
	"github.com/ucarion/urlpath"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type DolusExpectationEngine struct {
	cueExpectationsFiles []string
	expectations         map[expectation.Route][]expectation.DolusExpectation
	rawCueExpectations   []dolus.Expectation
	ResponseSchemas      map[expectation.Route]dstruct.DynamicStructModifier
	GenerationConfig     generator.GenerationConfig
	expectationRoutes    []string
}

var _ ExpectationEngine = &DolusExpectationEngine{}

func NewDolusExpectationEngine(
	generationConfig generator.GenerationConfig,
) (dolusExpectationEngine *DolusExpectationEngine) {
	dolusExpectationEngine = &DolusExpectationEngine{}
	dolusExpectationEngine.expectations = make(map[expectation.Route][]expectation.DolusExpectation)
	dolusExpectationEngine.ResponseSchemas = make(
		map[expectation.Route]dstruct.DynamicStructModifier,
	)
	dolusExpectationEngine.GenerationConfig = generationConfig

	return
}

func (e *DolusExpectationEngine) AddExpectation(
	expect expectation.DolusExpectation,
	validateExpectationSchema bool,
) error {
	// TODO check if exception overrides another one i.e has same request matcher

	if validateExpectationSchema {
		if err := e.validateExpectationSchema(expect); err != nil {
			return err
		}
	}
	route := expect.Request.Route
	e.expectations[route] = append(e.expectations[route], expect)

	if expect.RawCueExpectation != nil {
		e.rawCueExpectations = append(e.rawCueExpectations, *expect.RawCueExpectation)
	}

	return nil
}

func (e *DolusExpectationEngine) validateExpectationSchema(exp expectation.DolusExpectation) error {
	matchingResponseSchema, err := e.getMatchingResponseSchemaForRoute(exp.Request.Route)
	if err != nil {
		return err
	}
	expectationFieldErrors := validateExpectationResponseSchema(
		exp.Response.Body,
		matchingResponseSchema,
	)

	if len(expectationFieldErrors) > 0 {
		return expectation.ExpectationError{
			DolusExpectation:       exp,
			ExpectationFieldErrors: expectationFieldErrors,
		}
	}
	return nil
}

func (e *DolusExpectationEngine) AddResponseSchemaForRoute(
	route expectation.Route,
	responseSchema dstruct.DynamicStructModifier,
) error {
	if e.ResponseSchemas[route] != nil {
		return fmt.Errorf("response schema already exists for... ")
	}

	e.ResponseSchemas[route] = responseSchema
	return nil
}

func (e *DolusExpectationEngine) GetAllExpectations() map[expectation.Route][]expectation.DolusExpectation {
	return e.expectations
}

// TODO: either return []Route or change name to GetExpectationRoutePaths()
func (e *DolusExpectationEngine) GetExpectationRoutes() []string {
	return e.expectationRoutes
}

func (e *DolusExpectationEngine) GetExpectation(
	route expectation.Route,
) []expectation.DolusExpectation {
	return e.expectations[route]
}

func getParsedUrl(urlString string) (*url.URL, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return parsedURL, nil
}

func (e *DolusExpectationEngine) getMatchingResponseSchemaForRoute(
	route expectation.Route,
) (dstruct.DynamicStructModifier, error) {
	var matchingSchemas []dstruct.DynamicStructModifier

	parsedURL, err := getParsedUrl(route.Path)
	if err != nil {
		return nil, err
	}

	// expectationQueryParameters := parsedURL.Query()
	expectationPath := parsedURL.Path
	for schemaRoute, responseSchema := range e.ResponseSchemas {
		if schemaRoute.Match(route) {
			matchingSchemas = append(matchingSchemas, responseSchema)
		}
		// if pms.Method != k.Method || pms.Status != k.Status {

		// continue
		// }

		schemaPath := urlpath.New(schemaRoute.Path)
		_, ok := schemaPath.Match(expectationPath)
		if ok {
			matchingSchemas = append(matchingSchemas, responseSchema)
		}
	}

	if len(matchingSchemas) > 1 {
		return nil, fmt.Errorf("too many schemas match %+v", matchingSchemas)
	}
	if len(matchingSchemas) == 0 {
		return nil, fmt.Errorf("no matching schema found for path=%s", ucarionUrlPath)
	}

	return matchingSchemas[0], nil
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
							Err:       fmt.Errorf("invalid value '%s' for enum field: valid types are: \n%s", expectationFields[field].GetValue(), enumValues),
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

func (e *DolusExpectationEngine) getExpectaionsForRequest(
	path string,
	request *http.Request,
) []expectation.DolusExpectation {
	fmt.Println("COOL: ", path, request.URL.Path, request.RequestURI)
	// check for exact matches (with query parameters)
	if expectations := e.GetExpectation(expectation.Route{
		Path:   request.RequestURI,
		Method: request.Method,
	}); len(expectations) > 0 {
		return expectations
	}
	// get partial match if no exact match (with no query parameters)
	if expectations := e.GetExpectation(expectation.Route{
		Path:   request.URL.Path,
		Method: request.Method,
	}); len(expectations) > 0 {
		return expectations
	}

	return e.GetExpectation(expectation.Route{
		Path:   path,
		Method: request.Method,
	})
}

func getUcarionUrlPath(path string) string {
	p := strings.ReplaceAll(path, "{", ":")
	return strings.ReplaceAll(p, "}", "")
}

func (e *DolusExpectationEngine) GetResponseForRequest(
	path string,
	request *http.Request,
) (*expectation.DolusResponse, error) {
	// fmt.Println(path, method, request.URL.Path)
	// fmt.Println(len(expectations))
	expectations := e.getExpectaionsForRequest(path, request)

	if len(expectations) == 0 {
		return nil, fmt.Errorf("no expectation found for path and HTTP method")
	}

	// TODO find the right expectation depending on request matchers and priority
	currentExpectation := expectations[0]
	for _, v := range expectations {
		if request.URL.Path == v.Request.OpenApiPath {
			currentExpectation = v
			return &currentExpectation.Response, nil
		}
		if v.Priority > currentExpectation.Priority {
			currentExpectation = v
		}
	}

	return &currentExpectation.Response, nil
}

func (e *DolusExpectationEngine) GetRawCueExpectations() dolus.Expectations {
	return dolus.Expectations{
		Expectations: e.rawCueExpectations,
	}
}