package expectation

import (
	"net/http"

	"github.com/MartinSimango/dstruct"
	"github.com/ucarion/urlpath"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
)

type Route struct {
	Path   string
	Method string
}

func (r Route) Match(path string) bool {
	schemaPath := urlpath.New(r.Path)
	_, ok := schemaPath.Match(path)
	return ok
}

type DolusResponse struct {
	Body    dstruct.GeneratedStruct
	Status  int
	Headers http.Header
}

type DolusRequest struct {
	Route
	Body    any
	Headers http.Header
	// TODO: add cookies
}

type DolusExpectation struct {
	Priority          int
	Response          DolusResponse
	Request           DolusRequest
	RawCueExpectation *dolus.Expectation
}
