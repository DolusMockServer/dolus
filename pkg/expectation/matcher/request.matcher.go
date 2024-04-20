package matcher

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation/models"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type RequestMatcher struct {
	SimpleMatcher[models.Request]
	HeaderMatcher           SimpleMatcher[StringArrayMatcher]
	CookieMatcher           SimpleMatcher[CookieMatcher]
	RequestParameterMatcher SimpleMatcher[RequestParameterMatcher]
}

var _ Matcher[models.Request] = &RequestMatcher{}

func NewRequestMatcher(value models.Request) *RequestMatcher {
	return &RequestMatcher{
		SimpleMatcher: SimpleMatcher[models.Request]{
			MatchExpression: "eq",
			Value:           value,
		},
	}

}

func (rm RequestMatcher) MatchWithHttpRequest(request *http.Request, requestParameters schema.RequestParameters) bool {
	// match header
	return (rm.Value.Parameters == nil || rm.matchRequestParameters(requestParameters)) && rm.matchHeaders(request.Header) && rm.matchCookies(request)
}

func (rm RequestMatcher) Matches(value models.Request) bool {

	return true
}

// func (r *RequestMatcher) match(req *http.Request, reqParam schema.RequestParameters) bool {

// 	return (r.Parameters == nil || r.Parameters.Match(reqParam)) && r.matchHeaders(req.Header) &&
// 		r.matchCookies(req)
// }

func (r *RequestMatcher) matchHeaders(headers map[string][]string) bool {
	for name, value := range r.Value.Headers {
		requestHeaderValues := headers[name]
		if len(requestHeaderValues) == 0 {
			logger.Log.Debugf("No match for expectation! Header '%s' not present", name)
			return false
		}
		headerMatcher := value.(*StringArrayMatcher)
		if !(headerMatcher).Matches(requestHeaderValues) {
			logger.Log.Debugf("No match for expectation! Header '%s' with value %v does not match %v", name, headerMatcher.Value, requestHeaderValues)
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
		if !(cookieMatcherValue).MatchesWithRequestCookie(requestCookieValue) {
			logger.Log.Debugf("No match for expectation! Cookie '%s' with value %+v does not match %+v", cookieMatcherValue.Value.Name, cookieMatcherValue, *requestCookieValue)
			return false
		}
	}

	return true

}

func (r *RequestMatcher) matchRequestParameters(rp schema.RequestParameters) bool {
	return NewRequestParameterMatcher(*r.Value.Parameters).MatchWith(rp)
}
