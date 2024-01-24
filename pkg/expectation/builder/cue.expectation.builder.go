package builder

import (
	"fmt"
	"reflect"

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
	for _, instance := range *spec {
		expectations = append(expectations, ceb.buildExpectationFromCueInstance(instance)...)
	}
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
	for e.Next() {
		var cueExpectation expectation.Expectation
		err := e.Value().Decode(&cueExpectation)
		if err != nil {
			logger.Log.Error("Error decoding expectation: ", err)
			continue
		}
		if err := decodeMatcherFields(&cueExpectation); err != nil {
			logger.Log.Error("Error marshalling fields into matcher: ", err)
			continue
		}
		cueExpectation.Response.GeneratedBody = dstruct.NewGeneratedStructWithConfig(
			schema.SchemaFromAny(cueExpectation.Response.Body),
			&ceb.fieldGenerator,
		)
		expectations = append(expectations, cueExpectation)

	}
	return
}

// decodeMatcherFields decodes the matcher fields in the cueExpectation.
func decodeMatcherFields(cueExpectation *expectation.Expectation) (err error) {
	if err = createMatcherForRequestField[[]string](cueExpectation.Request.Headers); err != nil {
		return
	}
	if cueExpectation.Request.Parameters != nil {
		if err = createMatcherForRequestField[string](cueExpectation.Request.Parameters.Path); err != nil {
			return
		}
		if err = createMatcherForRequestField[[]string](cueExpectation.Request.Parameters.Query); err != nil {
			return
		}
	}
	return nil
}

// decodeMatcherField decodes the matcher fields in the map.
func createMatcherForRequestField[T any](m map[string]any) error {
	for name, value := range m {
		if v, err := createMatcher[T](value); err != nil {
			return fmt.Errorf("error with field '%s'= %v: %w", name, value, err)
		} else {
			m[name] = *v
		}

	}
	return nil
}

func createMatcher[T any](mapValue interface{}) (*expectation.Matcher[T], error) {
	switch field := mapValue.(type) {
	case map[string]interface{}:
		return createMatcherFromMatcherMap[T](field)
	case []interface{}:
		return createMatcherFromArrayValue[T](field, "eq")
	case interface{}:
		return createMatcherFromSingleValue[T](field, "eq")
	default:
		return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", mapValue, field)
	}
}

func createMatcherFromMatcherMap[T any](matcherMap map[string]interface{}) (*expectation.Matcher[T], error) {
	match := matcherMap["match"].(string)

	switch value := matcherMap["value"].(type) {
	case []interface{}:
		return createMatcherFromArrayValue[T](value, match)
	case interface{}:
		return createMatcherFromSingleValue[T](value, match)
	case nil:
		return &expectation.Matcher[T]{Match: match, Value: nil}, nil
	default:
		return nil, fmt.Errorf("could not marshal into Matcher: %v. Unsupported type %v", matcherMap, value)
	}
}

func createMatcherFromArrayValue[T any](arrayValue []interface{}, match string) (*expectation.Matcher[T], error) {
	var value []string
	for _, item := range arrayValue {
		value = append(value, fmt.Sprintf("%v", item))
	}

	return convertToMatcher[T](value, match)
}

func createMatcherFromSingleValue[T any](singleValue interface{}, match string) (*expectation.Matcher[T], error) {
	var value any
	switch any(*new(T)).(type) {
	case []string:
		value = []string{fmt.Sprintf("%v", singleValue)}
	default:
		value = fmt.Sprintf("%v", singleValue)
	}

	return convertToMatcher[T](value, match)
}

func convertToMatcher[T any](value any, match string) (*expectation.Matcher[T], error) {
	v, ok := value.(T)
	if !ok {
		return nil, fmt.Errorf("%v (%s) cannot be converted into %s",
			value,
			reflect.TypeOf(value).String(),
			reflect.TypeOf(v))
	}

	return &expectation.Matcher[T]{Match: match, Value: &v}, nil
}
