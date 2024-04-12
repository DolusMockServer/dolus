package expectation

import (
	"net/http"
	"net/url"

	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type Request struct {
	Path       string             `json:"path"`
	Method     string             `json:"method"`
	Body       any                `json:"body"`
	Parameters *RequestParameters `json:"params"`
	Headers    map[string]any     `json:"headers"`
	Cookies    []any              `json:"cookies"`
}

func (r *Request) Match(req *http.Request, reqParam schema.RequestParameters) bool {
	return (r.Parameters == nil || r.Parameters.Match(reqParam)) && r.matchHeaders(req.Header) &&
		r.matchCookies(req)
}

func (r *Request) matchHeaders(headers http.Header) bool {
	for name, value := range r.Headers {
		requestHeaderValues := headers[name]
		if len(requestHeaderValues) == 0 {
			logger.Log.Debugf("No match for expectation! Header '%s' not present", name)
			return false
		}
		if !(value.(StringArrayMatcher)).Matches(requestHeaderValues) {
			logger.Log.Debugf("No match for expectation! Header '%s' with value %v does not match %v", name, (value.(SimpleMatcher[[]string])).Value, requestHeaderValues)
			return false
		}
	}
	return true
}

func (r *Request) matchCookies(request *http.Request) bool {
	for _, cookie := range r.Cookies {
		cookieMatcherValue := cookie.(CookieMatcher)
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

func (r *Request) Route() schema.Route {
	return schema.Route{
		Path:   r.Path,
		Method: r.Method,
	}
}

func (r *Request) RouteWithParsedPath() (*schema.Route, error) {
	parsedURL, err := url.Parse(r.Path) // get rid of any query parameters
	if err != nil {
		return nil, err
	}
	return &schema.Route{
		Path:   schema.PathFromOpenApiPath(parsedURL.Path),
		Method: r.Method,
	}, nil
}
