package matcher

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation/models"
)

type CookieMatcher struct {
	CueMatcher[models.Cookie]
}

var _ Matcher[models.Cookie, http.Cookie] = &CookieMatcher{}

func NewCookieMatcher(value models.Cookie, matchType string) *CookieMatcher {
	return &CookieMatcher{
		CueMatcher: CueMatcher[models.Cookie]{
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
