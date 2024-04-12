package expectation

import "reflect"

type Matcher interface {
	Matches(value any) bool
	GetValue() any
}

type MatcherBuilder interface {
	Create(field map[string]any) (Matcher, error)
	CreateFromArrayValue(value []any, matchExpr string) (Matcher, error)
	CreateFromSingleValue(value any, matchExpr string) (Matcher, error)
}

var _ Matcher = &SimpleMatcher[string]{}

type SimpleMatcher[T any] struct {
	MatchExpression string `json:"match"`
	Value           T      `json:"value"` // TODO think about why this needs to be a pointer
}

// Matches implements Matcher.
func (m SimpleMatcher[T]) Matches(value any) bool {
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(m.Value, value)
		// case "has":
		// 	return m.Value != nil
	}
	return false
}

func (m SimpleMatcher[T]) GetValue() any {
	return m.Value
}
