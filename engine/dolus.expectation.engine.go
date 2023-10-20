package engine

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/MartinSimango/dolus-expectations/pkg/dolus"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/ucarion/urlpath"
)

type DolusExpectationEngine struct {
	cueExpectationsFiles []string
	expectations         map[expectation.PathMethod][]expectation.Expectation
	rawCueExpectations   []dolus.Expectation
	ResponseSchemas      map[expectation.PathMethodStatus]dstruct.DynamicStructModifier
	GenerationConfig     generator.GenerationConfig
}

var _ ExpectationEngine = &DolusExpectationEngine{}

func NewDolusExpectationEngine(generationConfig generator.GenerationConfig) (dolusExpectationEngine *DolusExpectationEngine) {
	dolusExpectationEngine = &DolusExpectationEngine{}
	dolusExpectationEngine.expectations = make(map[expectation.PathMethod][]expectation.Expectation)
	dolusExpectationEngine.ResponseSchemas = make(map[expectation.PathMethodStatus]dstruct.DynamicStructModifier)
	dolusExpectationEngine.GenerationConfig = generationConfig

	return
}

func (e *DolusExpectationEngine) AddExpectation(expect expectation.Expectation, validateExpectationSchema bool) error {
	// TODO check if exception overrides another one i.e has same request matcher
	pathMethod := expectation.PathMethod{
		Path:   expect.Request.Path,
		Method: expect.Request.Method,
	}
	if validateExpectationSchema {
		if err := e.validateExpectationSchema(expect); err != nil {
			return err
		}
	}

	e.expectations[pathMethod] = append(e.expectations[pathMethod], expect)

	fmt.Println("ADDING", pathMethod)
	if expect.RawCueExpectation != nil {
		e.rawCueExpectations = append(e.rawCueExpectations, *expect.RawCueExpectation)
	}

	return nil
}

func (e *DolusExpectationEngine) validateExpectationSchema(exp expectation.Expectation) error {
	matchingResponseSchema, err := e.getMatchingResponseSchemaForPathMethodStatus(expectation.PathMethodStatusExpectation(exp))
	if err != nil {
		return err
	}
	expectationFieldErrors := validateExpectationResponseSchema(exp.Response.Body, matchingResponseSchema)

	if len(expectationFieldErrors) > 0 {
		return expectation.ExpectationError{
			Expectation:            exp,
			ExpectationFieldErrors: expectationFieldErrors,
		}
	}
	return nil

}

func (e *DolusExpectationEngine) AddResponseSchemaForPathMethodStatus(pathMethodStatus expectation.PathMethodStatus, schema dstruct.DynamicStructModifier) error {

	if e.ResponseSchemas[pathMethodStatus] != nil {
		return fmt.Errorf("response schema already exists for... ")
	}

	e.ResponseSchemas[pathMethodStatus] = schema
	return nil
}

func (e *DolusExpectationEngine) GetExpectations() map[expectation.PathMethod][]expectation.Expectation {
	return e.expectations
}

func (e *DolusExpectationEngine) GetExpectationForPathMethod(pathMethod expectation.PathMethod) []expectation.Expectation {
	return e.expectations[pathMethod]
}

func (e *DolusExpectationEngine) getMatchingResponseSchemaForPathMethodStatus(pms expectation.PathMethodStatus) (dstruct.DynamicStructModifier, error) {
	var matchingSchemas []dstruct.DynamicStructModifier
	for k, v := range e.ResponseSchemas {
		if pms.Method != k.Method || pms.Status != k.Status {
			continue
		}
		schemaPath := urlpath.New(k.Path)
		_, ok := schemaPath.Match(pms.Path)
		if ok {
			matchingSchemas = append(matchingSchemas, v)
		}
	}

	if len(matchingSchemas) > 1 {
		return nil, fmt.Errorf("too many schemas match %+v", matchingSchemas)

	}
	if len(matchingSchemas) == 0 {
		return nil, fmt.Errorf("no matching schema found")

	}
	return matchingSchemas[0], nil

}

func addFieldDoesNotExistError(fieldName string, expectationFieldErrors []expectation.ExpectationFieldError) []expectation.ExpectationFieldError {
	fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
	return append(expectationFieldErrors, expectation.ExpectationFieldError{
		FieldName: fieldName,
		Err:       fmt.Errorf("field does not exist in the schema"),
	})
}

func addIncompatibleTypesError(fieldName string, expectationFieldErrors []expectation.ExpectationFieldError, schemaType reflect.Type, expectationType reflect.Type) []expectation.ExpectationFieldError {
	fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
	return append(expectationFieldErrors, expectation.ExpectationFieldError{
		FieldName: fieldName,
		Err:       fmt.Errorf("incompatible types. '%s' field is defined as type '%s' in schema but in expectation is defined as type '%s' ", fieldName, schemaType, expectationType),
	})
}

func validateExpectationResponseSchema(expect dstruct.DynamicStructModifier, schema dstruct.DynamicStructModifier) (expectationFieldErrors []expectation.ExpectationFieldError) {
	// e := dstruct.New(expectation.GetSchema())
	expectationFields := expect.GetFields()
	for field, value := range schema.GetFields() {
		schemaFieldType := value.GetType()
		expectationFieldType := expectationFields[field].GetType()

		if expectationFields[field].GetFieldName() == "" && value.GetTag("required") == "true" {
			expectationFieldErrors = addFieldDoesNotExistError(field, expectationFieldErrors)
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
				//validate enum type
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

func (e *DolusExpectationEngine) getExpectations(path, method string, request *http.Request) []expectation.Expectation {
	// check for exact matches
	if expectations := e.expectations[expectation.PathMethod{
		Path:   request.URL.Path,
		Method: method,
	}]; len(expectations) > 0 {
		return expectations
	}

	return e.expectations[expectation.PathMethod{
		Path:   path,
		Method: method,
	}]

}

func (e *DolusExpectationEngine) GetResponseForRequest(path, method string, request *http.Request) (*expectation.Response, error) {

	// fmt.Println(path, method, request.URL.Path)
	// fmt.Println(len(expectations))
	expectations := e.getExpectations(path, method, request)
	if len(expectations) == 0 {
		return nil, fmt.Errorf("no expectation found for path and HTTP method")
	}
	// TODO find the right expectation depending on request matchers and priority
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

func (e *DolusExpectationEngine) GetRawCueExpectations() dolus.Expectations {
	return dolus.Expectations{
		Expectations: e.rawCueExpectations,
	}
}
