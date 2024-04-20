package matcher

import (
	"fmt"
	"reflect"
)

type MatcherBuilder[T any, M any] interface {
	Create(field map[string]any) (Matcher[T, M], error)
	CreateFromArrayValue(value []any, matchExpr string) (Matcher[T, M], error)
	CreateFromSingleValue(value any, matchExpr string) (Matcher[T, M], error)
}

func ConvertMapKeysToMatchers[T any, M any](builder MatcherBuilder[T, M], mapValue map[string]any) (err error) {

	for k, v := range mapValue {
		if mapValue[k], err = convertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert map field to matcher: %w", err)
		}
	}
	return nil

}

func ConvertArrayFieldsToMatchers[T any, M any](builder MatcherBuilder[T, M], arrayValue []any) (err error) {
	for i, v := range arrayValue {
		if arrayValue[i], err = convertToMatcher(v, builder); err != nil {
			return fmt.Errorf("failed to convert array field to matcher: %w", err)
		}
	}
	return nil
}

func convertToMatcher[T any, M any](v any, builder MatcherBuilder[T, M]) (Matcher[T, M], error) {
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
