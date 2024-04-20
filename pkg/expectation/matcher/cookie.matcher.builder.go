package matcher

import (
	"encoding/json"
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type CookieMatcherBuilder struct{}

var _ MatcherBuilder[expectation.Cookie, http.Cookie] = &CookieMatcherBuilder{}

func (b CookieMatcherBuilder) Create(field map[string]any) (Matcher[expectation.Cookie, http.Cookie], error) {

	data, _ := json.Marshal(field) // this should never fail as cue validated it
	matchExpr := "eq"
	var v expectation.Cookie
	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil

}
func (b CookieMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher[expectation.Cookie, http.Cookie], error) {
	panic("Cannot create a CookieMatcher from an array value")
}

func (b CookieMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher[expectation.Cookie, http.Cookie], error) {
	var v expectation.Cookie
	data, _ := json.Marshal(value) // this should never fail as cue validated it

	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil
}
