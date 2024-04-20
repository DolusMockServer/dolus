package matcher

import (
	"reflect"
)

type StringArrayMatcher struct {
	SimpleMatcher[[]string]
}

var _ Matcher[[]string] = &StringArrayMatcher{}

func NewStringArrayMatcher(value []string, matchType string) *StringArrayMatcher {
	return &StringArrayMatcher{
		SimpleMatcher: SimpleMatcher[[]string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}
}

func (m StringArrayMatcher) Matches(value []string) bool {
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(m.Value, value)
	case "has":
		return m.Value != nil && len(m.Value) > 0
	}
	return false
}
