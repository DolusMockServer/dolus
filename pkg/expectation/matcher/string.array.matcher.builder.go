package matcher

import (
	"encoding/json"
	"fmt"
)

type StringArrayMatcherBuilder struct{}

var _ MatcherBuilder[[]string, []string] = &StringArrayMatcherBuilder{}

func (b StringArrayMatcherBuilder) Create(field map[string]any) (Matcher[[]string, []string], error) {

	data, _ := json.Marshal(field["value"]) // this should never fail as cue validated it
	matchExpr := field["match"].(string)
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		var vs string
		if json.Unmarshal(data, &vs) == nil {
			stringArray := []string{vs}
			return NewStringArrayMatcher(&stringArray, matchExpr), nil
		}
		return nil, err
	}
	return NewStringArrayMatcher(&v, matchExpr), nil

}

func (b StringArrayMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher[[]string, []string], error) {
	var stringArray []string
	for _, v := range value {
		stringArray = append(stringArray, fmt.Sprintf("%v", v))
	}
	return NewStringArrayMatcher(&stringArray, matchExpr), nil
}

func (b StringArrayMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher[[]string, []string], error) {
	stringArray := []string{fmt.Sprintf("%v", value)}
	return NewStringArrayMatcher(&stringArray, matchExpr), nil
}
