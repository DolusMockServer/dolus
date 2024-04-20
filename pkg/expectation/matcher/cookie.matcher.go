package matcher

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type CookieMatcher struct {
	expectation.CueMatcher[expectation.Cookie]
}

var _ Matcher[expectation.Cookie, http.Cookie] = &CookieMatcher{}

func NewCookieMatcher(value expectation.Cookie, matchType string) *CookieMatcher {
	return &CookieMatcher{
		CueMatcher: expectation.CueMatcher[expectation.Cookie]{
			MatchExpression: matchType,
			Value:           &value,
		},
	}
}

func (m CookieMatcher) Matches(cookie *http.Cookie) bool {
	switch m.MatchExpression {
	case "eq":
		return cookie.Path == m.Value.Path && cookie.Value == m.Value.Value
	case "has":
		return m.Value != nil
	}
	return false
}
