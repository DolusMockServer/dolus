package matcher

import (
	"github.com/DolusMockServer/dolus/pkg/expectation/models"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type RequestParameterMatcher struct {
	SimpleMatcher[models.RequestParameters]
}

var _ Matcher[models.RequestParameters] = &RequestParameterMatcher{}

func NewRequestParameterMatcher(value models.RequestParameters) *RequestParameterMatcher {
	return &RequestParameterMatcher{
		SimpleMatcher: SimpleMatcher[models.RequestParameters]{
			MatchExpression: "eq",
			Value:           value,
		},
	}

}

func (r *RequestParameterMatcher) Match(rp models.RequestParameters) bool {
	return true
}

func (r *RequestParameterMatcher) MatchWith(rp schema.RequestParameters) bool {
	return r.matchPathParams(rp.PathParams) && r.matchQueryParams(rp.QueryParams)
}

// TODO see if we can make queryParams map[string]*StringArrayMatcher
func (r *RequestParameterMatcher) matchQueryParams(queryParams map[string][]string) bool {
	// query parameters are already validated so we can just check if the values match
	for name, value := range r.Value.Query {
		matcher := value.(*StringArrayMatcher)
		matches := matcher.Matches(queryParams[name])
		if !matches {
			logger.Log.Debugf("No match for expectation! Query parameter '%s' with value %v does not match %v", name, matcher.Value, queryParams[name])
			return false
		}
	}
	return true

}

func (r *RequestParameterMatcher) matchPathParams(pathParams map[string]string) bool {
	// query parameters are already validated so we can just check if the values match
	for name, value := range r.Value.Path {
		matcher := value.(*StringMatcher)
		matches := matcher.Matches(pathParams[name])
		if !matches {
			logger.Log.Debugf("No match for expectation! Path parameter '%s' with value %v does not match %v", name, matcher.Value, pathParams[name])
			return false
		}
	}
	return true

}
