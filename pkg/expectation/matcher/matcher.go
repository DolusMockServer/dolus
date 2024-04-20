package matcher

import "reflect"

type Matcher[T any] interface {
	Matches(value T) bool
	GetValue() T
}

var _ Matcher[string] = &SimpleMatcher[string]{}

type SimpleMatcher[T any] struct {
	MatchExpression string `json:"match"`
	Value           T      `json:"value"` // TODO think about why this needs to be a pointer
}

// Matches implements Matcher.
func (m SimpleMatcher[T]) Matches(value T) bool {
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(m.Value, value)
		// case "has":
		// 	return m.Value != nil
	}
	return false
}

func (m SimpleMatcher[T]) GetValue() T {
	return m.Value
}
