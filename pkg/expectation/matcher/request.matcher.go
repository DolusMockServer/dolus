package matcher

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type RequestMatcher struct {
	Value                   expectation.Request
	RequestParameterMatcher *RequestParameterMatcher
}

func NewRequestMatcher(value expectation.Request) *RequestMatcher {
	return &RequestMatcher{
		Value:                   value,
		RequestParameterMatcher: NewRequestParameterMatcher(value.Parameters),
	}

}

func (rm RequestMatcher) Matches(request *http.Request, requestParameters *schema.RequestParameters) bool {

	return (rm.Value.Parameters == nil || rm.RequestParameterMatcher.Matches(requestParameters)) && rm.matchHeaders(request.Header) && rm.matchCookies(request)
}

func (r *RequestMatcher) matchHeaders(headers map[string][]string) bool {
	for name, value := range r.Value.Headers {
		requestHeaderValues := headers[name]
		if len(requestHeaderValues) == 0 {
			logger.Log.Debugf("No match for expectation! Header '%s' not present", name)
			return false
		}
		headerMatcher := value.(*StringArrayMatcher)
		if !(headerMatcher).Matches(&requestHeaderValues) {
			logger.Log.Debugf("No match for expectation! Header '%s' with value %v does not match %v", name, *headerMatcher.Value, requestHeaderValues)
			return false
		}
	}
	return true
}

func (r *RequestMatcher) matchCookies(request *http.Request) bool {
	for _, cookie := range r.Value.Cookies {
		cookieMatcherValue := cookie.(*CookieMatcher)
		requestCookieValue, err := request.Cookie(cookieMatcherValue.Value.Name)
		if err != nil {
			logger.Log.Debugf("No match for expectation! Cookie '%s' not present", cookieMatcherValue.Value.Name)
			return false
		}
		if !(cookieMatcherValue).Matches(requestCookieValue) {
			logger.Log.Debugf("No match for expectation! Cookie '%s' with value %+v does not match %+v", cookieMatcherValue.Value.Name, cookieMatcherValue, *requestCookieValue)
			return false
		}
	}

	return true

}
