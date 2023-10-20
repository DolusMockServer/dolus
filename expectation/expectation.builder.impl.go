package expectation

import (
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dolus-expectations/pkg/dolus"
	"github.com/MartinSimango/dolus/core"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"
)

type ExpectationBuilderImpl struct {
	loader         Loader[OpenAPISpecLoadType]
	fieldGenerator generator.Generator
}

var _ ExpectationBuilder = &ExpectationBuilderImpl{}

func NewExpectationBuilderImpl(fieldGenerator generator.Generator) *ExpectationBuilderImpl {
	return &ExpectationBuilderImpl{
		fieldGenerator: fieldGenerator,
	}
}

func (eb *ExpectationBuilderImpl) BuildExpectationsFromCueLoader(loader Loader[CueExpectationLoadType]) ([]Expectation, error) {
	if s, err := loader.load(); err != nil {
		return nil, err
	} else {
		return eb.buildExpectationsFromCueLoadType(s), nil
	}
}

func (eb *ExpectationBuilderImpl) BuildExpectationsFromOpenApiSpecLoader(loader Loader[OpenAPISpecLoadType]) ([]Expectation, error) {

	if s, err := loader.load(); err != nil {
		return nil, err
	} else {
		return eb.buildExpectationsFromOpenApiSpec(s), nil
	}
}

func getRealPath(path string) string {
	p := strings.ReplaceAll(path, "{", ":")
	return strings.ReplaceAll(p, "}", "")
}

func (eb *ExpectationBuilderImpl) buildExpectationsFromOpenApiSpec(spec *OpenAPISpecLoadType) (expectations []Expectation) {
	for path := range spec.Paths {
		for method, operation := range spec.Paths[path].Operations() {
			p := getRealPath(path)
			for code, ref := range operation.Responses {
				if p != "/store/order/:orderId" || code != "200" {
					continue
				}
				// if p != "/" || code != "200" {
				// 	continue
				// }
				fmt.Println(p, code)
				responseSchema := core.NewResponseSchemaFromOpenApi3Ref(p, method, code, ref, "application/json")

				// engine must store for each path method code then check that

				status, _ := strconv.Atoi(code)

				expectations = append(expectations, Expectation{
					Priority: 0,
					Request: Request{
						Path:   getRealPath(path),
						Method: method,
					},
					Response: Response{
						Body:   dstruct.NewGeneratedStructWithConfig(responseSchema.Schema.GetSchema(), &eb.fieldGenerator),
						Status: status,
					},
				})

			}
		}

	}
	return
}

func (eb *ExpectationBuilderImpl) buildExpectationsFromCueLoadType(spec *CueExpectationLoadType) (expectations []Expectation) {
	for _, instance := range *spec {
		expectations = append(expectations, eb.buildExpectationFromCueInstance(instance)...)
	}
	return
}

func (eb *ExpectationBuilderImpl) buildExpectationFromCueInstance(instance cue.Value) (expectations []Expectation) {
	e, err := instance.Value().LookupPath(cue.ParsePath("expectations")).List()

	if err != nil {
		fmt.Printf("error with expectation in file %s: %s \n", instance.Pos().Filename(), err)
		return
	}
	for e.Next() {
		var cueExpectation dolus.Expectation
		err := e.Value().Decode(&cueExpectation)
		if err != nil {
			fmt.Println("Error decoding expectation: ", err)
			continue
		}
		status := strconv.Itoa(cueExpectation.Response.Status)
		// TODO schema
		r := core.NewResponseSchemaFromAny(cueExpectation.Request.Path, cueExpectation.Request.Method, status, cueExpectation.Response.Body)

		expectations = append(expectations, Expectation{
			RawCueExpectation: &cueExpectation,
			Priority:          cueExpectation.Priority,
			Response: Response{
				Body:   dstruct.NewGeneratedStructWithConfig(r.GetSchema(), &eb.fieldGenerator),
				Status: cueExpectation.Response.Status,
			},
			Request: Request{
				Path:   cueExpectation.Request.Path,
				Method: cueExpectation.Request.Method,
			},
		})

	}
	return

}
