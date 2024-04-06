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
	if err = createMatcherForRequestField[expectation.StringArrayMatcher](cueExpectation.Request.Headers); err != nil {
		return
	}
	if cueExpectation.Request.Parameters != nil {
		if err = createMatcherForRequestField[expectation.StringMatcher](cueExpectation.Request.Parameters.Path); err != nil {
			return
		}
		if err = createMatcherForRequestField[expectation.StringArrayMatcher](cueExpectation.Request.Parameters.Query); err != nil {
			return
		}
	}
	if err = createMatcherForRequestFieldArray[expectation.CookieMatcher](cueExpectation.Request.Cookies); err != nil {
		return
	}
	return nil
}

func createStringArrayMatcherFromMap(_map map[string]any) error {
	for key, value := range _map {
		// 1 determine if it is a matcher or not
		switch field := value.(type) {
		case map[string]interface{}: // matcher type
			fieldValue, err := json.Marshal(field["value"]) // this should never fail as cue validated it
			if err != nil {
				return err
			}
			if err := json.Unmarshal(fieldValue, &value); err != nil {
				_map[key] = expectation.NewStringArrayMatcher([]string{fmt.Sprintf("%v", value)}, field["match"].(string))
			}

		case []interface{}:
			var v []string
			for _, item := range field {
				v = append(v, fmt.Sprintf("%v", item))
			}
			_map[key] = expectation.NewStringArrayMatcher(v, "eq")

		case interface{}:
			_map[key] = expectation.NewStringArrayMatcher([]string{fmt.Sprintf("%v", value)}, "eq")

		default:
			return fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", value, field)
		}

		// 2 if not matcher, convert to matcher
	}

	return nil
}

func createStringMatcherFromMap(_map map[string]any) error {
	for key, value := range _map {
		// 1 determine if it is a matcher or not
		switch field := value.(type) {
		case map[string]interface{}: // matcher type

		case interface{}:
			_map[key] = expectation.NewStringArrayMatcher([]string{fmt.Sprintf("%v", value)}, "eq")

		default:
			return fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", value, field)
		}

		// 2 if not matcher, convert to matcher
	}

	return nil
}

// decodeMatcherField decodes the matcher fields in the map.
func createMatcherForRequestField[T expectation.MatcherType](m map[string]any) error {
	for name, value := range m {
		if v, err := createMatcherFromMapValue[T](value); err != nil {
			return fmt.Errorf("error with field '%s'= %v: %w", name, value, err)
		} else {
			m[name] = v
		}

	}
	return nil

}

// cookies: [
// 	expectation.#Cookie & {name: "cookie1", value: "As", path: "/p"},
// 	expectation.#CookieMatcher & { match}

// ]

func createMatcherForRequestFieldArray[T expectation.MatcherType](m []any) error {
	for i, value := range m {
		if v, err := createMatcherFromMapValue[T](value); err != nil {
			return fmt.Errorf("error with field '%d'= %v: %w", i, value, err)
		} else {
			m[i] = v
			// 2.
		}

	}
	return nil
}

func createMatcherFromMapValue[T expectation.MatcherType](mapValue interface{}) (expectation.MatcherType, error) {

	switch field := mapValue.(type) {
	case map[string]interface{}: // matcher type
		var value T

		fieldValue, err := json.Marshal(field["value"]) // this should never fail as cue validated it
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(fieldValue, &value); err != nil {
			return createMatcherFromMatcherMap[T](field)
		}

		return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", mapValue, field)

	case []interface{}:

		return createMatcherFromArrayValue[T](field, "eq")
	case interface{}:
		return createMatcherFromSingleValue[T](field, "eq")

	default:
		return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", mapValue, field)
	}
}

func createMatcherFromMatcherMap[T expectation.MatcherType](matcherMap map[string]interface{}) (expectation.MatcherType, error) {

	v, ok := any(matcherMap).(expectation.Matcher[T])
	if ok {
		fmt.Printf("YESS : %v\n", v)
	} else {
		fmt.Printf("NOO: %v\n", any(matcherMap))
	}

	match := matcherMap["match"].(string)

	switch value := matcherMap["value"].(type) {
	case []interface{}:
		return createMatcherFromArrayValue[T](value, match)
	case interface{}:
		return createMatcherFromSingleValue[T](value, match)
	case nil:
		return nil, nil
		// return &expectation.Matcher[T]{MatchExpression: match, Value: nil}, nil
	default:
		return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", matcherMap, value)
	}
}

func createMatcherFromArrayValue[T expectation.MatcherType](arrayValue []interface{}, match string) (expectation.MatcherType, error) {

	var value any

	switch any(*new(T)).(type) {
	case expectation.StringArrayMatcher:
		var v []string
		for _, item := range arrayValue {
			value = append(v, fmt.Sprintf("%v", item))
		}
		value = v
	case expectation.CookieMatcher:
		var v []expectation.Cookie
		for _, item := range arrayValue {
			value = append(v, item.(expectation.Cookie))
		}
		value = v

	}
	return convertToMatcher[T](value, match)
}

func createMatcherFromSingleValue[T expectation.MatcherType](singleValue interface{}, match string) (expectation.MatcherType, error) {
	var value any
	switch any(*new(T)).(type) {
	case expectation.StringArrayMatcher:
		value = []string{fmt.Sprintf("%v", singleValue)}
	case expectation.CookieMatcher:
		value = singleValue
	default:
		value = fmt.Sprintf("%v", singleValue)
	}

	return convertToMatcher[T](value, match)
}

func convertToMatcher[T expectation.MatcherType](value any, match string) (expectation.MatcherType, error) {
	switch any(*new(T)).(type) {
	case expectation.StringMatcher:
		v, ok := value.(string)
		if !ok {
			return nil, cannotConvertError(value, v)
		}
		return expectation.NewStringMatcher(v, match), nil

	case expectation.StringArrayMatcher:
		v, ok := value.([]string)
		if !ok {
			return nil, cannotConvertError(value, v)
		}
		return expectation.NewStringArrayMatcher(v, match), nil
	case expectation.CookieMatcher:
		v, ok := value.(expectation.Cookie)
		if !ok {
			return nil, cannotConvertError(value, v)
		}
		return expectation.NewCookieMatcher(v, match), nil
	}

	return nil, fmt.Errorf("unknown matcher type '%v'", reflect.TypeOf(*new(T)).String())
}

func cannotConvertError(value any, v any) error {
	return fmt.Errorf("%v of type '%s' cannot be converted into '%s'",
		value,
		reflect.TypeOf(value).String(),
		reflect.TypeOf(v))
}
