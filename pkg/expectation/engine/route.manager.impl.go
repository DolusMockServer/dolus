package engine

import (
	"github.com/DolusMockServer/dolus/pkg/schema"
	"github.com/MartinSimango/dstruct"
)

type RouteManagerImpl struct {
	routeProperties map[schema.Route]*schema.RouteProperty
}

var _ RouteManager = &RouteManagerImpl{}

func NewRouteManager() *RouteManagerImpl {
	return &RouteManagerImpl{
		routeProperties: make(map[schema.Route]*schema.RouteProperty),
	}
}

// AddRoute implements RouteManager.
func (r *RouteManagerImpl) AddRoute(route schema.Route, routeProperty *schema.RouteProperty) {
	r.routeProperties[route] = routeProperty
}

// RemoveRoute implements RouteManager.
func (r *RouteManagerImpl) RemoveRoute(route schema.Route) {
	r.routeProperties[route] = nil
}

// GetRequestSchema implements RouteManager.
func (r *RouteManagerImpl) GetRequestSchema(route schema.Route) dstruct.DynamicStructModifier {
	return r.routeProperties[route].RequestSchema
}

// GetResponseSchema implements RouteManager.
func (r *RouteManagerImpl) GetResponseSchema(route schema.Route) dstruct.DynamicStructModifier {
	return r.routeProperties[route].ResponseSchema
}

// GetRoutes implements RouteManager.
func (r *RouteManagerImpl) GetRoutes() []schema.Route {
	routes := make([]schema.Route, 0, len(r.routeProperties))
	for route := range r.routeProperties {
		routes = append(routes, route)
	}
	return routes
}

// GetRouteProperty implements RouteManager.
func (r *RouteManagerImpl) GetRouteProperty(route schema.Route) *schema.RouteProperty {
	return r.routeProperties[route]
}
