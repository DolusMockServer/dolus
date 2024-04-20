package matcher

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type RequestParameterMatcher struct {
	expectation.CueMatcher[expectation.RequestParameters]
}

var _ Matcher[expectation.RequestParameters, schema.RequestParameters] = &RequestParameterMatcher{}

func NewRequestParameterMatcher(value *expectation.RequestParameters) *RequestParameterMatcher {
	return &RequestParameterMatcher{
		CueMatcher: expectation.CueMatcher[expectation.RequestParameters]{
			MatchExpression: "eq",
			Value:           value,
		},
	}

}

func (r *RequestParameterMatcher) Matches(rp *schema.RequestParameters) bool {
	return r.matchPathParams(rp.PathParams) && r.matchQueryParams(rp.QueryParams)
}

// TODO see if we can make queryParams map[string]*StringArrayMatcher
func (r *RequestParameterMatcher) matchQueryParams(queryParams map[string][]string) bool {
	// query parameters are already validated so we can just check if the values match
	for name, value := range r.Value.Query {
		matcher := value.(*StringArrayMatcher)
		param := queryParams[name]
		matches := matcher.Matches(&param)
		if !matches {
			logger.Log.Debugf("No match for expectation! Query parameter '%s' with value %v does not match %v", name, *matcher.Value, queryParams[name])
			return false
		}
	}
	return true

}

func (r *RequestParameterMatcher) matchPathParams(pathParams map[string]string) bool {
	// query parameters are already validated so we can just check if the values match
	for name, value := range r.Value.Path {
		matcher := value.(*StringMatcher)
		param := pathParams[name]
		matches := matcher.Matches(&param)
		if !matches {
			logger.Log.Debugf("No match for expectation! Path parameter '%s' with value %v does not match %v", name, *matcher.Value, pathParams[name])
			return false
		}
	}
	return true

}
