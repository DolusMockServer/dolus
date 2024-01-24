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
	Cookies    *[]Cookie          `json:"cookies"`
}

func (r *Request) Match(req *http.Request, reqParam schema.RequestParameters) bool {
	return (r.Parameters == nil || r.Parameters.Match(reqParam)) && r.matchHeaders(req.Header)
}

func (r *Request) matchHeaders(headers http.Header) bool {
	for name, value := range r.Headers {
		requestHeaderValues := headers[name]
		if len(requestHeaderValues) == 0 {
			logger.Log.Debugf("No match for expectation! Header '%s' not found", name)
			return false
		}
		if !(value.(Matcher[[]string])).Matches(&requestHeaderValues) {
			logger.Log.Debugf("No match for expectation! Header '%s' with value %v does not match %v", name, *(value.(Matcher[[]string])).Value, requestHeaderValues)
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
