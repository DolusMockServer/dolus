package engine

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type ExpectationEngine interface {
	AddExpectation(expectation expectation.Expectation,
		validateExpectationSchema bool,
		expectationType expectation.ExpectationType) error
	AddResponseSchemaForRoute(
		route schema.Route,
		responseSchema dstruct.DynamicStructModifier,
	) error
	GetAllExpectations() map[schema.Route][]expectation.Expectation
	GetExpectation(route schema.Route) []expectation.Expectation
	GetResponseForRequest(
		request *http.Request,
		reqParam schema.RequestParameters,
		path string,
	) (*expectation.Response, error)
	GetCueExpectations() expectation.Expectations
	GetExpectationRoutes() []schema.Route
	SetRouteProperties(routeProperties schema.RouteProperties)
	AddRoute(route schema.Route) error
	GetRoutes() []schema.Route
}
