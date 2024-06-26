package matcher

import (
	"reflect"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type StringArrayMatcher struct {
	expectation.CueMatcher[[]string]
}

var _ Matcher[[]string, []string] = &StringArrayMatcher{}

func NewStringArrayMatcher(value *[]string, matchType string) *StringArrayMatcher {
	return &StringArrayMatcher{
		CueMatcher: expectation.CueMatcher[[]string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}
}

func (m StringArrayMatcher) Matches(value *[]string) bool {
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(*m.Value, *value)
	case "has":
		return m.Value != nil && len(*m.Value) > 0
	}
	return false
}
