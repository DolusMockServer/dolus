package matcher

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation/models"
)

type CookieMatcher struct {
	SimpleMatcher[models.Cookie]
}

var _ Matcher[models.Cookie] = &CookieMatcher{}

func NewCookieMatcher(value models.Cookie, matchType string) *CookieMatcher {
	return &CookieMatcher{
		SimpleMatcher: SimpleMatcher[models.Cookie]{
			MatchExpression: matchType,
			Value:           value,
		},
	}
}

func (m CookieMatcher) Matches(cookie models.Cookie) bool {

	switch m.MatchExpression {
	case "eq":
		return cookie.Path == m.Value.Path && cookie.Value == m.Value.Value
	case "has":
		return true
	}
	return false
}

func (m CookieMatcher) MatchesWithRequestCookie(cookie *http.Cookie) bool {
	switch m.MatchExpression {
	case "eq":
		return cookie.Path == m.Value.Path && cookie.Value == m.Value.Value
	case "has":
		return true
	}
	return false
}
