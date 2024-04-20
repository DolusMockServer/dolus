package builder

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
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
			if err := decodeMatcherFields(&cueExpectation); err != nil {
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

// func createHeaderMatchRule(cueExpectation *models.Expectation) (err error) {
// 	for k, v := range cueExpectation.Request.Headers {
// 		var m matcher.Matcher[[]string]
// 		if m, err = ConvertToMatcher(v, matcher.StringArrayMatcherBuilder{}); err != nil {
// 			return fmt.Errorf("failed to convert map field to matcher: %w", err)
// 		}
// 		cueExpectation.MatchRules.Headers[k] = models.Rule[[]string]{
// 			MatchType: m.(*matcher.StringArrayMatcher).MatchExpression,
// 			Value:     m.GetValue(),
// 		}
// 	}
// 	return nil
// }

// decodeMatcherFields decodes the matcher fields in the cueExpectation.
func decodeMatcherFields(cueExpectation *models.Expectation) (err error) {

	if err = ConvertMapKeysToMatchers(matcher.StringArrayMatcherBuilder{}, cueExpectation.Request.Headers); err != nil {
		return
	}
	if cueExpectation.Request.Parameters != nil {
		if err = ConvertMapKeysToMatchers(matcher.StringMatcherBuilder{}, cueExpectation.Request.Parameters.Path); err != nil {
			return
		}
		if err = ConvertMapKeysToMatchers(matcher.StringArrayMatcherBuilder{}, cueExpectation.Request.Parameters.Query); err != nil {
			return
		}
	}
	if err = ConvertArrayFieldsToMatchers(matcher.CookieMatcherBuilder{}, cueExpectation.Request.Cookies); err != nil {
		return
	}
	return nil
}

func ConvertMapKeysToMatchers[T any](builder matcher.MatcherBuilder[T], mapValue map[string]any) (err error) {

	for k, v := range mapValue {
		if mapValue[k], err = ConvertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert map field to matcher: %w", err)
		}
	}
	return nil

}

func ConvertArrayFieldsToMatchers[T any](builder matcher.MatcherBuilder[T], arrayValue []any) (err error) {
	for i, v := range arrayValue {
		if arrayValue[i], err = ConvertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert array field to matcher: %w", err)
		}
	}
	return nil
}

func ConvertToMatcher[T any](v any, builder matcher.MatcherBuilder[T]) (matcher.Matcher[T], error) {
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
		expectation.Request.Parameters.Query[k] = matcher.NewStringArrayMatcher(value, "eq")
	}
	return nil
}
