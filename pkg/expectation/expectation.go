package expectation

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
)

type Route struct {
	Path   string
	Method string
}

func (r Route) Match(route Route) bool {
	panic("Not implemented")
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
