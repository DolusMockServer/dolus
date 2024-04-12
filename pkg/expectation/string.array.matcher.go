package expectation

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type StringArrayMatcher struct {
	SimpleMatcher[[]string]
}

var _ Matcher = &StringArrayMatcher{}

func NewStringArrayMatcher(value []string, matchType string) *StringArrayMatcher {
	return &StringArrayMatcher{
		SimpleMatcher: SimpleMatcher[[]string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}
}

func (m StringArrayMatcher) Matches(value any) bool {
	stringArray := value.([]string)
	switch m.MatchExpression {
	case "eq":
		return reflect.DeepEqual(m.Value, stringArray)
	case "has":
		return m.Value != nil && len(m.Value) > 0
	}
	return false
}

type StringArrayMatcherBuilder struct{}

var _ MatcherBuilder = &StringArrayMatcherBuilder{}

func (b StringArrayMatcherBuilder) Create(field map[string]any) (Matcher, error) {

	data, _ := json.Marshal(field["value"]) // this should never fail as cue validated it
	matchExpr := field["match"].(string)
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		var vs string
		if json.Unmarshal(data, &vs) == nil {
			return NewStringArrayMatcher([]string{vs}, matchExpr), nil
		}
		return nil, err
	}
	return NewStringArrayMatcher(v, matchExpr), nil

}

func (b StringArrayMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher, error) {
	var stringArray []string
	for _, v := range value {
		stringArray = append(stringArray, fmt.Sprintf("%v", v))
	}
	return NewStringArrayMatcher(stringArray, matchExpr), nil
}

func (b StringArrayMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher, error) {
	return NewStringArrayMatcher([]string{fmt.Sprintf("%v", value)}, matchExpr), nil
}
