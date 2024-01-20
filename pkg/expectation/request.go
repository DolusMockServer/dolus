package expectation

import (
	"net/url"

	"github.com/DolusMockServer/dolus/pkg/schema"
)

type Request struct {
	Path       string               `json:"path"`
	Method     string               `json:"method"`
	Body       any                  `json:"body"`
	Parameters *RequestParameters   `json:"params"`
	Headers    *map[string][]string `json:"headers"`
	Cookies    *[]Cookie            `json:"cookies"`
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
