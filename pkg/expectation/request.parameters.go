package expectation

import (
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

const (
	QUERY_PARAM string = "Query"
	PATH_PARAM  string = "Path"
)

type RequestParameters struct {
	Path  map[string]any `json:"path"`
	Query map[string]any `json:"query"`
}

func (r *RequestParameters) Match(rp schema.RequestParameters) bool {
	return matchParams(PATH_PARAM, r.Path, rp.PathParams) && matchParams(QUERY_PARAM, r.Query, rp.QueryParams)
}

func matchParams[T any](pathType string, params map[string]any, values map[string]T) bool {
	for name, value := range params {
		v := values[name]
		if !(value.(Matcher[T])).Matches(&v) {
			logger.Log.Debugf("No match for expectation! %s parameter '%s' with value %v does not match %v", pathType, name, *(value.(Matcher[T]).Value), v)
			return false
		}
	}
	return true
}
