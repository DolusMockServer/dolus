package builder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"

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
func (ceb *CueExpectationBuilder) BuildExpectations() (*Output, error) {
	if s, err := ceb.loader.Load(); err != nil {
		return nil, err
	} else {
		return &Output{
			Expectations: ceb.buildExpectationsFromCueLoadType(s),
		}, nil
	}
}

func (ceb *CueExpectationBuilder) buildExpectationsFromCueLoadType(
	spec *loader.CueExpectationLoadType,
) (expectations []expectation.Expectation) {
	t := time.Now()
	for _, instance := range *spec {
		expectations = append(expectations, ceb.buildExpectationFromCueInstance(instance)...)
	}
	fmt.Printf("Time to build %d expectations: %v\n", len(expectations), time.Since(t))
	return
}

func (ceb *CueExpectationBuilder) buildExpectationFromCueInstance(
	instance cue.Value,
) (expectations []expectation.Expectation) {
	e, err := instance.Value().LookupPath(cue.ParsePath("expectations")).List()
	if err != nil {
		fmt.Printf("error with expectation in file %s: %s \n", instance.Pos().Filename(), err)
		return
	}
	var wg sync.WaitGroup
	for e.Next() {
		wg.Add(1)
		go func(cueValue cue.Value) {
			defer wg.Done()
			var cueExpectation expectation.Expectation

			err := cueValue.Decode(&cueExpectation)
			if err != nil {
				logger.Log.Error("Error decoding expectation: ", err)
				// continue
				return
			}
			if err := decodeMatcherFields(&cueExpectation); err != nil {
				logger.Log.Error("Error marshalling fields into matcher: ", err)
				// continue
				return
			}
			a, _ := json.Marshal(cueExpectation.Request.Headers)
			fmt.Printf("AFTER: %v\n", string(a))

			b, _ := json.Marshal(cueExpectation.Request.Cookies)
			fmt.Printf("AFTER: %v\n", string(b))

			cueExpectation.Response.GeneratedBody = dstruct.NewGeneratedStructWithConfig(
				schema.SchemaFromAny(cueExpectation.Response.Body),
				&ceb.fieldGenerator,
			)
			expectations = append(expectations, cueExpectation)
		}(e.Value())

	}
	wg.Wait()
	return
}

// decodeMatcherFields decodes the matcher fields in the cueExpectation.
func decodeMatcherFields(cueExpectation *expectation.Expectation) (err error) {

	if err = ConvertMapKeysToMatchers(expectation.StringArrayMatcherBuilder{}, cueExpectation.Request.Headers); err != nil {
		return
	}
	if cueExpectation.Request.Parameters != nil {
		if err = ConvertMapKeysToMatchers(expectation.StringMatcherBuilder{}, cueExpectation.Request.Parameters.Path); err != nil {
			return
		}
		if err = ConvertMapKeysToMatchers(expectation.StringArrayMatcherBuilder{}, cueExpectation.Request.Parameters.Query); err != nil {
			return
		}
	}
	if err = ConvertArrayFieldsToMatchers(expectation.CookieMatcherBuilder{}, cueExpectation.Request.Cookies); err != nil {
		return
	}
	return nil
}

func ConvertMapKeysToMatchers(builder expectation.MatcherBuilder, mapValue map[string]any) (err error) {

	for k, v := range mapValue {
		if mapValue[k], err = ConvertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert map field to matcher: %w", err)
		}
	}
	return nil

}

func ConvertArrayFieldsToMatchers(builder expectation.MatcherBuilder, arrayValue []any) (err error) {
	for i, v := range arrayValue {
		if arrayValue[i], err = ConvertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert array field to matcher: %w", err)
		}
	}
	return nil
}

func ConvertToMatcher(v any, builder expectation.MatcherBuilder) (expectation.Matcher, error) {
	switch field := v.(type) {
	case map[string]interface{}:
		return builder.Create(field)
	case []interface{}:
		return builder.CreateFromArrayValue(field, "eq")
	case interface{}:
		return builder.CreateFromSingleValue(field, "eq")

	}
	return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", v, reflect.TypeOf(v))
}
