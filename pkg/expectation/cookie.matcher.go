package expectation

import (
	"encoding/json"
)

type CookieMatcher struct {
	SimpleMatcher[Cookie]
}

func NewCookieMatcher(value Cookie, matchType string) *CookieMatcher {
	return &CookieMatcher{
		SimpleMatcher: SimpleMatcher[Cookie]{
			MatchExpression: matchType,
			Value:           value,
		},
	}
}

var _ Matcher = &CookieMatcher{}

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

type CookieMatcherBuilder struct{}

var _ MatcherBuilder = &CookieMatcherBuilder{}

func (b CookieMatcherBuilder) Create(field map[string]any) (Matcher, error) {

	data, _ := json.Marshal(field) // this should never fail as cue validated it
	matchExpr := "eq"
	var v Cookie
	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil

}
func (b CookieMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher, error) {
	panic("Cannot create a CookieMatcher from an array value")
}

func (b CookieMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher, error) {
	var v Cookie
	data, _ := json.Marshal(value) // this should never fail as cue validated it

	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil
}
