package engine

import (
	"fmt"
	"net/http"

	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/ucarion/urlpath"
)

type DolusExpectationEngine struct {
	cueExpectationsFiles []string
	expectations         map[expectation.PathMethod][]expectation.Expectation
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

	return nil
}

func (e *DolusExpectationEngine) validateExpectationSchema(exp expectation.Expectation) error {
	matchingResponseSchema, err := e.getMatchingResponseSchemaForPathMethodStatus(expectation.PathMethodStatusExpectation(exp))
	if err != nil {
		return fmt.Errorf("error with expectation: %s", err)
	}

	if !doesResponseSchemaMatch(exp.Response.Body, matchingResponseSchema) {
		return fmt.Errorf("schema does not match")
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

func doesResponseSchemaMatch(expectation dstruct.DynamicStructModifier, schema dstruct.DynamicStructModifier) bool {
	// e := dstruct.New(expectation.GetSchema())

	// err := dstruct.DoSchemasMatch(e, dstruct.New(schema.GetSchema()))
	// e.Print()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return false
	// }

	// fmt.Println(reflect.responseSchema1.GetSchema(), responseSchema2.GetSchema())
	return true
}

func (e *DolusExpectationEngine) GetResponseForRequest(path, method string, request *http.Request) (*expectation.Response, error) {
	expectations := e.expectations[expectation.PathMethod{
		Path:   request.URL.Path,
		Method: method,
	}]
	// fmt.Println(path, method, request.URL.Path)
	// fmt.Println(len(expectations))
	if len(expectations) == 0 {

		expectations = e.expectations[expectation.PathMethod{
			Path:   path,
			Method: method,
		}]

		if len(expectations) == 0 {
			return nil, fmt.Errorf("no expectation found for path and HTTP method")
		}
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
