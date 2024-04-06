package expectation

import (
	"reflect"
)

type MatcherType interface {
	Matches(value any) bool
	GetValue() any
}

type Matcher[T any] struct {
	MatchExpression string `json:"match"`
	Value           *T     `json:"value"`
}

func (m Matcher[T]) GetValue() any {
	return m.Value
}

type StringMatcher struct {
	Matcher[string]
}

func NewStringMatcher(value, matchType string) *StringMatcher {
	return &StringMatcher{
		Matcher: Matcher[string]{
			MatchExpression: matchType,
			Value:           &value,
		},
	}

}

type StringArrayMatcher struct {
	Matcher[[]string]
}

func NewStringArrayMatcher(value []string, matchType string) *StringArrayMatcher {
	return &StringArrayMatcher{
		Matcher: Matcher[[]string]{
			MatchExpression: matchType,
			Value:           &value,
		},
	}
}

type CookieMatcher struct {
	Matcher[Cookie]
}

func NewCookieMatcher(value Cookie, matchType string) *CookieMatcher {
	return &CookieMatcher{
		Matcher: Matcher[Cookie]{
			MatchExpression: matchType,
			Value:           &value,
		},
	}
}

var _ MatcherType = &StringMatcher{}
var _ MatcherType = &StringArrayMatcher{}
var _ MatcherType = &CookieMatcher{}

func (m StringMatcher) Matches(value any) bool {
	switch m.MatchExpression {
	case "eq":
		return *m.Value == value.(string)
	case "has":
		return true
	}
	return false
}

func (m StringArrayMatcher) Matches(value any) bool {
	stringArray := value.([]string)
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(m.Value, stringArray)
	case "has":
		return m.Value != nil && len(*m.Value) > 0
	}
	return false
}

func (m CookieMatcher) Matches(value any) bool {
	cookie := value.(*Cookie)
	switch m.MatchExpression {
	case "eq":
		return cookie.Path == m.Value.Path && cookie.Value == m.Value.Value
	case "has":
		return true
	}
	return false
}
