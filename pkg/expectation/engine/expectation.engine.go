package engine

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type ExpectationEngine interface {
	// AddExpectation adds an expectation to the engine
	AddExpectation(expectation expectation.Expectation,
		validateExpectationSchema bool,
		expectationType expectation.ExpectationType) error

	// AddRoute registers a new route with the engine. If the route already exists, it will be overwritten.
	AddRoute(route schema.Route) error

	// AddResponseSchemaForRoute adds a schema for a route's response body. If the route has not been registered, it will return an error.
	AddResponseSchema(
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
	GetRoutes() []schema.Route
}
