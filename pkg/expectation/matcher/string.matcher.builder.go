package matcher

import (
	"encoding/json"
	"fmt"
)

type StringMatcherBuilder struct{}

var _ MatcherBuilder[string, string] = &StringMatcherBuilder{}

func (b StringMatcherBuilder) Create(field map[string]any) (Matcher[string, string], error) {

	data, _ := json.Marshal(field["value"]) // this should never fail as cue validated it
	matchExpr := field["match"].(string)
	var v string
	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewStringMatcher(&v, matchExpr), nil

}

func (b StringMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher[string, string], error) {
	v := fmt.Sprintf("%v", value)
	return NewStringMatcher(&v, matchExpr), nil
}

func (b StringMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher[string, string], error) {
	v := fmt.Sprintf("%v", value)
	return NewStringMatcher(&v, matchExpr), nil
}
