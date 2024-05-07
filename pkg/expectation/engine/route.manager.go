package engine

import (
	"github.com/DolusMockServer/dolus/pkg/schema"
	"github.com/MartinSimango/dstruct"
)

type RouteManager interface {
	// AddRoute registers a new route with the route manager. If the route already exists, it will be overwritten.
	AddRoute(route schema.Route,

		routeProperty *schema.RouteProperty)
	// RemoveRoute removes a route from the manager.
	RemoveRoute(route schema.Route)

	// GetRoutes returns all the routes registered with the manager.
	GetRoutes() []schema.Route

	GetRouteProperties() map[schema.Route]*schema.RouteProperty

	// GetRouteProperty returns the route property for the given route.
	GetRouteProperty(route schema.Route) *schema.RouteProperty

	// GetRequestSchema returns the request schema for the given route.
	GetRequestSchema(route schema.Route) dstruct.DynamicStructModifier

	// GetResponseSchema returns the response schema for the given route.
	GetResponseSchema(route schema.Route) dstruct.DynamicStructModifier
}
