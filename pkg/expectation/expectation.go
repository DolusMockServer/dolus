package expectation

import (
	"strings"

	"github.com/DolusMockServer/dolus/pkg/expectation/cue"
	"github.com/MartinSimango/dstruct"
	"github.com/ucarion/urlpath"
)

type Route struct {
	Path      string
	Operation string
}

func (r Route) Match(path string) bool {
	schemaPath := urlpath.New(r.Path)
	m, ok := schemaPath.Match(path)
	if ok {
		// path Params must have perfect matches i.e /order/:id should not match /order/:ida
		for k, v := range m.Params {
			if strings.HasPrefix(v, ":") {
				if strings.ToLower(k) != strings.ToLower(v[1:]) {
					return false
				}
			}
		}
	}
	return ok
}

type DolusResponse struct {
	Body dstruct.GeneratedStruct
}

type DolusRequest struct {
	cue.Request
	Route
}
