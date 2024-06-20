package builder

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/loader"
	"github.com/DolusMockServer/dolus/pkg/expectation/matcher"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type CueExpectationBuilder struct {
	loader                    loader.Loader[loader.CueExpectationLoadType]
	fieldGenerator            generator.Generator
	cookieMatcherBuilder      matcher.MatcherBuilder[expectation.Cookie, http.Cookie]
	stringArrayMatcherBuilder matcher.MatcherBuilder[[]string, []string]
	stringMatcherBuilder      matcher.MatcherBuilder[string, string]
}

// check that we implement the interface
var _ ExpectationBuilder = &CueExpectationBuilder{}

func NewCueExpectationBuilder(
	cueExpectationFiles []string,
	fieldGenerator generator.Generator,
) *CueExpectationBuilder {
	return &CueExpectationBuilder{
		loader:                    loader.NewCueExpectationLoader(cueExpectationFiles),
		fieldGenerator:            fieldGenerator,
		cookieMatcherBuilder:      matcher.CookieMatcherBuilder{},
		stringArrayMatcherBuilder: matcher.StringArrayMatcherBuilder{},
		stringMatcherBuilder:      matcher.StringMatcherBuilder{},
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
			if err := ceb.decodeMatcherFields(&cueExpectation); err != nil {
				logger.Log.Error("Error marshalling fields into matcher: ", err)
				// continue
				return
			}

			cueExpectation.Response.GeneratedBody = dstruct.NewGeneratedStructWithConfig(
				schema.SchemaFromAny(cueExpectation.Response.Body),
				&ceb.fieldGenerator,
			)
			cueExpectation.ExpectationType = expectation.Custom
			expectations = append(expectations, cueExpectation)
		}(e.Value())

	}
	wg.Wait()
	return
}

func (ceb *CueExpectationBuilder) decodeMatcherFields(cueExpectation *expectation.Expectation) (err error) {

	if err = matcher.ConvertMapKeysToMatchers(ceb.stringArrayMatcherBuilder, cueExpectation.Request.Headers); err != nil {
		return
	}
	if cueExpectation.Request.Parameters != nil {
		if err = matcher.ConvertMapKeysToMatchers(ceb.stringMatcherBuilder, cueExpectation.Request.Parameters.Path); err != nil {
			return
		}
		if err = matcher.ConvertMapKeysToMatchers(ceb.stringArrayMatcherBuilder, cueExpectation.Request.Parameters.Query); err != nil {
			return
		}
	}
	if err = matcher.ConvertArrayFieldsToMatchers(ceb.cookieMatcherBuilder, cueExpectation.Request.Cookies); err != nil {
		return
	}
	return nil
}
