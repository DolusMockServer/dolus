package builder

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/loader"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type CueExpectationBuilder struct {
	loader         loader.Loader[loader.CueExpectationLoadType]
	fieldGenerator generator.Generator
}

// check that we implement the interface
var _ ExpectationBuilder = &CueExpectationBuilder{}

func NewCueExpectationBuilder(
	cueExpectationFiles []string,
	fieldGenerator generator.Generator,
) *CueExpectationBuilder {
	return &CueExpectationBuilder{
		loader:         loader.NewCueExpectationLoader(cueExpectationFiles),
		fieldGenerator: fieldGenerator,
	}
}

// BuildExpectations implements ExpectationBuilder.
func (ceb *CueExpectationBuilder) BuildExpectations() ([]expectation.DolusExpectation, error) {
	if s, err := ceb.loader.Load(); err != nil {
		return nil, err
	} else {
		return ceb.buildExpectationsFromCueLoadType(s), nil
	}
}

func (ceb *CueExpectationBuilder) buildExpectationsFromCueLoadType(
	spec *loader.CueExpectationLoadType,
) (expectations []expectation.DolusExpectation) {
	for _, instance := range *spec {
		expectations = append(expectations, ceb.buildExpectationFromCueInstance(instance)...)
	}
	return
}

func (ceb *CueExpectationBuilder) buildExpectationFromCueInstance(
	instance cue.Value,
) (expectations []expectation.DolusExpectation) {
	e, err := instance.Value().LookupPath(cue.ParsePath("expectations")).List()
	if err != nil {
		fmt.Printf("error with expectation in file %s: %s \n", instance.Pos().Filename(), err)
		return
	}
	for e.Next() {
		var cueExpectation dolus.Expectation
		err := e.Value().Decode(&cueExpectation)
		if err != nil {
			logger.Log.Error("Error decoding expectation: ", err)
			continue
		}

		expectations = append(expectations, expectation.DolusExpectation{
			CueExpectation: &cueExpectation,
			Priority:       cueExpectation.Priority,
			Request: expectation.DolusRequest{
				Route: expectation.Route{
					Path:   pathFromOpenApiPath(cueExpectation.Request.Path),
					Method: string(cueExpectation.Request.Method),
				},
				Body: cueExpectation.Request.Body,
			},
			Response: expectation.DolusResponse{
				Body: dstruct.NewGeneratedStructWithConfig(
					schema.SchemaFromAny(cueExpectation.Response.Body),
					&ceb.fieldGenerator,
				),
				Status: cueExpectation.Response.Status,
			},
		})

	}
	return
}
