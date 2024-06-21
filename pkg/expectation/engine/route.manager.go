package engine

import (
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type RouteManagerError interface {
	error
}

type RouteManager interface {
	// AddRoute registers a new route with the route manager. If the route already exists, it will be overwritten.
	AddRoute(route schema.Route,
		routeProperty schema.RouteProperty)
	// RemoveRoute removes a route from the manager.
	RemoveRoute(route schema.Route)

	// GetRoutes returns all the routes registered with the manager.
	GetRoutes() []schema.Route

	// GetRouteProperty returns the route property for the given route and an error if the route does not exist.
	GetRouteProperty(route schema.Route) (schema.RouteProperty, RouteManagerError)

	// GetRouteProperties returns all the route properties registered with the manager.
	GetRouteProperties() map[schema.Route]schema.RouteProperty
}
