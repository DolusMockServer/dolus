package engine

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type ExpectationEngine interface {
	AddExpectation(expectation expectation.DolusExpectation, validateExpectationSchema bool) error
	AddResponseSchemaForRoute(
		route expectation.Route,
		responseSchema dstruct.DynamicStructModifier,
	) error
	GetAllExpectations() map[expectation.Route][]expectation.DolusExpectation
	GetExpectation(route expectation.Route) []expectation.DolusExpectation
	GetResponseForRequest(path string, request *http.Request) (*expectation.DolusResponse, error)
	// TODO: rename in dolus-expecation package to CueExpectation
	GetCueExpectations() dolus.Expectations
	GetExpectationRoutes() []string
}
