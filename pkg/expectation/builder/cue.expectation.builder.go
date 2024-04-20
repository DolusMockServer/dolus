package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus/pkg/expectation/loader"
	"github.com/DolusMockServer/dolus/pkg/expectation/matcher"
	"github.com/DolusMockServer/dolus/pkg/expectation/models"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type CueExpectationBuilder struct {
	loader                    loader.Loader[loader.CueExpectationLoadType]
	fieldGenerator            generator.Generator
	cookieMatcherBuilder      matcher.MatcherBuilder[models.Cookie, http.Cookie]
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
) (expectations []models.Expectation) {
	t := time.Now()
	for _, instance := range *spec {
		expectations = append(expectations, ceb.buildExpectationFromCueInstance(instance)...)
	}
	fmt.Printf("Time to build %d expectations: %v\n", len(expectations), time.Since(t))
	return
}

func (ceb *CueExpectationBuilder) buildExpectationFromCueInstance(
	instance cue.Value,
) (expectations []models.Expectation) {
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
			var cueExpectation models.Expectation

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
			// add query parameters from path to expectation and overrides old query parameters from cue file
			addQueryParameters(&cueExpectation)

			a, _ := json.Marshal(cueExpectation.Request.Headers)
			fmt.Printf("AFTER: %v\n", string(a))

			b, _ := json.Marshal(cueExpectation.Request.Cookies)
			fmt.Printf("AFTER COOKIES: %v\n", string(b))

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

func (ceb *CueExpectationBuilder) decodeMatcherFields(cueExpectation *models.Expectation) (err error) {

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

func addQueryParameters(expectation *models.Expectation) error {
	parsedURL, err := url.Parse(expectation.Request.Path)
	if err != nil {
		return fmt.Errorf("failed to add query parameters for expectation with path '%s': %w", expectation.Request.Path, err)

	}
	queryParams := parsedURL.Query()
	if expectation.Request.Parameters == nil {
		expectation.Request.Parameters = &models.RequestParameters{}
	}
	if expectation.Request.Parameters.Query == nil {
		expectation.Request.Parameters.Query = make(map[string]any)
	}
	for k, v := range queryParams {
		value := v
		expectation.Request.Parameters.Query[k] = matcher.NewStringArrayMatcher(&value, "eq")
	}
	return nil
}
