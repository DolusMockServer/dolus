package engine

import (
	"fmt"
	"net/http"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/MartinSimango/dolus/core"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dolus/generator"
	"github.com/ucarion/urlpath"
)

const CueExpectationModuleRoot = "/Users/martinsimango/dolus/expectations"

type DolusExpectationEngine struct {
	cueExpectationsFiles []string
	expectations         map[expectation.PathMethod][]expectation.Expectation
	ResponseSchemas      map[expectation.PathMethodStatus]*core.ResponseSchema
	GenerationConfig     generator.GenerationConfig
}

var _ ExpectationEngine = &DolusExpectationEngine{}

func NewDolusExpectationEngine(generationConfig generator.GenerationConfig) (dolusExpectationEngine *DolusExpectationEngine) {
	dolusExpectationEngine = &DolusExpectationEngine{}
	dolusExpectationEngine.expectations = make(map[expectation.PathMethod][]expectation.Expectation)
	dolusExpectationEngine.ResponseSchemas = make(map[expectation.PathMethodStatus]*core.ResponseSchema)
	dolusExpectationEngine.GenerationConfig = generationConfig

	return
}

func (e *DolusExpectationEngine) AddResponseSchemaForPathMethod(responseSchema *core.ResponseSchema) error {
	pathMethodStatus := expectation.PathMethodStatus{
		PathMethod: expectation.PathMethod{
			Path:   responseSchema.Path,
			Method: responseSchema.Method,
		},
		Status: responseSchema.StatusCode,
	}
	if e.ResponseSchemas[pathMethodStatus] != nil {
		return fmt.Errorf("response schema already exists for... ")
	}
	e.ResponseSchemas[pathMethodStatus] = responseSchema
	return nil
}

func (e *DolusExpectationEngine) AddExpectationsFromFiles(files ...string) {
	// TODO get files extention to see which expectation files to load
	e.cueExpectationsFiles = append(e.cueExpectationsFiles, files...)
}

func (e *DolusExpectationEngine) Load() error {
	// TODO needs to check if expectations match schemas
	return e.loadCueExpectation()
}

func (e *DolusExpectationEngine) loadCueExpectation() error {
	ctx := cuecontext.New()
	entrypoints := e.cueExpectationsFiles

	bis := load.Instances(entrypoints, &load.Config{
		ModuleRoot: CueExpectationModuleRoot,
	})

	for _, bi := range bis {
		// check for errors on the  instance
		// these are typically parsing errors
		if bi.Err != nil {
			fmt.Println("Error during load:", bi.Err)
			continue
		}
		value := ctx.BuildInstance(bi)

		if value.Err() != nil {
			fmt.Println("Error during build:", value.Err())
			continue
		}

		// Validate the value
		err := value.Validate()
		if err != nil {
			fmt.Println("Error during validation:", err)
			continue
		}

		e.addExpectationFromCueValue(value)
	}
	return nil
}

func (e *DolusExpectationEngine) AddExpectation(pathMethod expectation.PathMethod, expectation expectation.Expectation) error {
	// TODO check if exception overrides another one i.e has same request matcher
	e.expectations[pathMethod] = append(e.expectations[pathMethod], expectation)
	return nil
}

func (e *DolusExpectationEngine) GetExpectations() map[expectation.PathMethod][]expectation.Expectation {
	return e.expectations
}

func (e *DolusExpectationEngine) GetExpectationForPathMethod(pathMethod expectation.PathMethod) []expectation.Expectation {
	return e.expectations[pathMethod]
}

func (e *DolusExpectationEngine) addExpectationFromCueValue(instance cue.Value) {

	expectations, _ := instance.Value().LookupPath(cue.ParsePath("expectations")).List()
	for expectations.Next() {

		var cueExpectation expectation.CueExpectation

		err := expectations.Value().Decode(&cueExpectation)
		if err != nil {
			fmt.Println("Error decoding expectation: ", err)
			continue
		}
		status := fmt.Sprintf("%d", cueExpectation.Response.Status)
		responseSchema := core.NewResponseSchemaFromAny(
			cueExpectation.Path,
			cueExpectation.Method,
			status,
			cueExpectation.Response.Body,
		)
		// TODO now check that scham matches responseSchames for path and status

		// pathMethod := expectation.PathMethod{
		// 	Path:   cueExpectation.Path,
		// 	Method: cueExpectation.Method,
		// }
		// pathMethodStatus := expectation.PathMethodStatus{
		// 	PathMethod: expectation.PathMethod{
		// 		Path:   responseSchema.Path,
		// 		Method: responseSchema.Method,
		// 	},
		// 	Status: responseSchema.StatusCode,
		// }
		// /store/order/1
		// /store/order/:orderId
		// /store/order
		sc, err := e.getMatchingPathSchema(responseSchema.Path, responseSchema.Method, responseSchema.StatusCode)
		if err != nil {
			fmt.Println("Error with expectation! ", err)
			continue
		}
		if doesResponseSchemaMatch(responseSchema, sc) {

			responseExample := core.NewExample(responseSchema, e.GenerationConfig)
			matchingPathMethod := expectation.PathMethod{
				Path:   sc.Path,
				Method: sc.Method,
			}

			e.expectations[matchingPathMethod] = append(e.expectations[matchingPathMethod], expectation.Expectation{
				Pririoty: cueExpectation.Pririoty,
				Response: expectation.Response{
					Body:   *responseExample,
					Status: cueExpectation.Response.Status,
				},
				Request: expectation.Request{
					Path:   responseSchema.Path,
					Method: responseSchema.Method,
				},
			})
		} else {
			fmt.Println("Schema does not match!!!")
		}

	}
}

func (e *DolusExpectationEngine) getMatchingPathSchema(path string, method string, statusCode string) (*core.ResponseSchema, error) {
	var matchingSchemas []*core.ResponseSchema
	for k, v := range e.ResponseSchemas {
		if method != k.Method || statusCode != k.Status {
			continue
		}
		schemaPath := urlpath.New(k.Path)
		_, ok := schemaPath.Match(path)
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

func doesResponseSchemaMatch(expectation *core.ResponseSchema, schema *core.ResponseSchema) bool {
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
		Path:   path,
		Method: method,
	}]
	fmt.Println(path, method, request.URL.Path)
	fmt.Println(len(expectations))
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
		if v.Pririoty > currentExpectation.Pririoty {
			currentExpectation = v
		}
	}

	return &currentExpectation.Response, nil
}
