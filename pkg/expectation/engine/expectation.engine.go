package engine

import (
	"net/http"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type ExpectationEngine interface {
	// AddExpectation adds an expectation to the engine
	AddExpectation(expectation expectation.Expectation,
		validateExpectationSchema bool) error

	// GetExpectations returns the expectations for the given expectation type. Expectations can be filtered by path, method and expectation type.
	GetExpectations(expectationType *expectation.ExpectationType, path *string, method *string) []expectation.Expectation

	// GetResponseForRequest returns the response for the given request
	GetResponseForRequest(
		request *http.Request,
		reqParams schema.RequestParameters,
		path string,
	) (*expectation.Response, error)

	// SetRouteManager sets the route manager for the engine
	SetRouteManager(routeManager RouteManager)

	// GetRouteManager returns the route manager for the engine
	GetRouteManager() RouteManager
}
