package schema

import (
	"reflect"
	"strings"

	"github.com/MartinSimango/dstruct"
	"github.com/ucarion/urlpath"
)

type Route struct {
	Path   string
	Method string
}

// Match returns true if the given path matches the route.
func (r Route) Match(path string) (map[string]string, bool) {
	schemaPath := urlpath.New(r.Path)
	m, ok := schemaPath.Match(path)
	if ok {
		// path Params should have perfect matches i.e /order/:id should not match /order/:ida
		for k, v := range m.Params {
			if strings.HasPrefix(v, ":") {
				if !strings.EqualFold(k, v[1:]) {
					return nil, false
				} else {
					return m.Params, true
				}
			}
		}
	}
	return m.Params, ok
}

type ParameterProperty struct {
	Required bool
	Type     reflect.Type
}

type ParameterProperties map[string]*ParameterProperty

type RequestParameterProperty struct {
	PathParameterProperties  ParameterProperties
	QueryParameterProperties ParameterProperties
}

type RouteProperty struct {
	RequestParameterProperty
	RequestSchema  dstruct.DynamicStructModifier
	ResponseSchema dstruct.DynamicStructModifier
}

// type RouteProperties map[Route]RouteProperty
