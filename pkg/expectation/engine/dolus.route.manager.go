package engine

import (
	"fmt"

	"github.com/DolusMockServer/dolus/pkg/schema"
)

// TODO: think about using caching to reduce constantly refectching routes and routeProperties

type DolusRouteManagerError struct {
	Message string
}

var _ RouteManagerError = &DolusRouteManagerError{}

func (rme *DolusRouteManagerError) Error() string {
	return rme.Message
}

type DolusRouteManager struct {
	routeProperties map[schema.Route]schema.RouteProperty
}

var _ RouteManager = &DolusRouteManager{}

func NewRouteManager() *DolusRouteManager {
	return &DolusRouteManager{
		routeProperties: make(map[schema.Route]schema.RouteProperty),
	}
}

// AddRoute implements RouteManager.
func (r *DolusRouteManager) AddRoute(route schema.Route, routeProperty schema.RouteProperty) {
	r.routeProperties[route] = routeProperty
}

// RemoveRoute implements RouteManager.
func (r *DolusRouteManager) RemoveRoute(route schema.Route) {
	delete(r.routeProperties, route)
}

// GetRoutes implements RouteManager.
func (r *DolusRouteManager) GetRoutes() []schema.Route {
	routes := make([]schema.Route, 0, len(r.routeProperties))
	for route := range r.routeProperties {
		routes = append(routes, route)
	}
	return routes
}

// GetRouteProperty implements RouteManager.
func (r *DolusRouteManager) GetRouteProperty(route schema.Route) (schema.RouteProperty, RouteManagerError) {
	if !r.doesRouteExist(route) {
		return schema.RouteProperty{}, &DolusRouteManagerError{
			Message: fmt.Sprintf("route %v does not exist", route),
		}
	}
	return r.routeProperties[route], nil
}

// GetRouteProperties implements RouteManager.
func (r *DolusRouteManager) GetRouteProperties() map[schema.Route]schema.RouteProperty {
	// returning a copy is less efficient but safer
	routePropertiesCopy := make(map[schema.Route]schema.RouteProperty, len(r.routeProperties))
	for route, prop := range r.routeProperties {
		routePropertiesCopy[route] = prop
	}
	return routePropertiesCopy

}

func (r *DolusRouteManager) doesRouteExist(route schema.Route) bool {
	_, ok := r.routeProperties[route]
	return ok
}
