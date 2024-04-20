package engine

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus/pkg/expectation/models"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type ExpectationEngine interface {
	AddExpectation(expectation models.Expectation,
		validateExpectationSchema bool,
		expectationType models.ExpectationType) error
	AddResponseSchemaForRoute(
		route schema.Route,
		responseSchema dstruct.DynamicStructModifier,
	) error
	GetAllExpectations() map[schema.Route][]models.Expectation
	GetExpectation(route schema.Route) []models.Expectation
	GetResponseForRequest(
		request *http.Request,
		reqParam schema.RequestParameters,
		path string,
	) (*models.Response, error)
	GetCueExpectations() models.Expectations
	GetExpectationRoutes() []schema.Route
	SetRouteProperties(routeProperties schema.RouteProperties)
}
