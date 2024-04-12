package expectation

import (
	"encoding/json"
	"fmt"
)

type StringMatcher struct {
	SimpleMatcher[string]
}

var _ Matcher = &StringMatcher{}

func NewStringMatcher(value, matchType string) *StringMatcher {
	return &StringMatcher{
		SimpleMatcher: SimpleMatcher[string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}

}

func (m StringMatcher) Matches(value any) bool {
	switch m.MatchExpression {
	case "eq":
		return m.Value == value.(string)
	case "has":
		return true
	}
	return false
}

type StringMatcherBuilder struct{}

var _ MatcherBuilder = &StringMatcherBuilder{}

func (b StringMatcherBuilder) Create(field map[string]any) (Matcher, error) {

	data, _ := json.Marshal(field["value"]) // this should never fail as cue validated it
	matchExpr := field["match"].(string)
	var v string
	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewStringMatcher(v, matchExpr), nil

}

func (b StringMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher, error) {
	return NewStringMatcher(fmt.Sprintf("%v", value), matchExpr), nil
}

func (b StringMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher, error) {
	return NewStringMatcher(fmt.Sprintf("%v", value), matchExpr), nil
}
