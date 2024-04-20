package matcher

import (
	"encoding/json"

	"github.com/DolusMockServer/dolus/pkg/expectation/models"
)

type CookieMatcherBuilder struct{}

var _ MatcherBuilder[models.Cookie] = &CookieMatcherBuilder{}

func (b CookieMatcherBuilder) Create(field map[string]any) (Matcher[models.Cookie], error) {

	data, _ := json.Marshal(field) // this should never fail as cue validated it
	matchExpr := "eq"
	var v models.Cookie
	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil

}
func (b CookieMatcherBuilder) CreateFromArrayValue(value []any, matchExpr string) (Matcher[models.Cookie], error) {
	panic("Cannot create a CookieMatcher from an array value")
}

func (b CookieMatcherBuilder) CreateFromSingleValue(value any, matchExpr string) (Matcher[models.Cookie], error) {
	var v models.Cookie
	data, _ := json.Marshal(value) // this should never fail as cue validated it

	if err := json.Unmarshal(data, &v); err != nil {

		return nil, err
	}
	return NewCookieMatcher(v, matchExpr), nil
}
