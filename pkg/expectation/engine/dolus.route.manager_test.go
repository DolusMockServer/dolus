//go:build unit
// +build unit

package engine

import (
	"testing"

	"github.com/MartinSimango/dstruct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DolusMockServer/dolus/pkg/schema"
)

type DolusRouteManagerTestSuite struct {
	suite.Suite
	routeManager RouteManager
}

func (suite *DolusRouteManagerTestSuite) SetupTest() {
	suite.routeManager = NewRouteManager()
}

func (suite *DolusRouteManagerTestSuite) TestAddRoute() {
	suite.T().Run("can add route", func(t *testing.T) {
		// Given
		suite.SetupTest()
		route := schema.Route{
			Path:   "/v1/test",
			Method: "GET",
		}
		routeProperty := schema.RouteProperty{}

		// When
		suite.routeManager.AddRoute(route, routeProperty)

		returnedRouteProperty, err := suite.routeManager.GetRouteProperty(route)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, routeProperty, returnedRouteProperty)
	})
}

func (suite *DolusRouteManagerTestSuite) TestRemoveRoute() {
	suite.T().Run("can remove route", func(t *testing.T) {
		// Given
		suite.SetupTest()
		route := schema.Route{
			Path:   "/v1/test",
			Method: "GET",
		}

		// When
		suite.routeManager.RemoveRoute(route)

		_, err := suite.routeManager.GetRouteProperty(route)

		// Then
		assert.Error(t, err)
	})
}

func (suite *DolusRouteManagerTestSuite) TestGetRoutes() {
	suite.T().Run("can get routes", func(t *testing.T) {
		// Given
		suite.SetupTest()
		routes := []schema.Route{
			{Path: "/v1/test", Method: "GET"},
			{Path: "/v1/test", Method: "POST"},
			{Path: "/v2/test", Method: "GET"},
		}
		for _, route := range routes {
			suite.routeManager.AddRoute(route, schema.RouteProperty{})
		}

		// When
		returnedRoutes := suite.routeManager.GetRoutes()

		// Then
		assert.ElementsMatch(t, routes, returnedRoutes)

	})

}

func (suite *DolusRouteManagerTestSuite) TestGetRouteProperty() {

	suite.T().Run("can get route property", func(t *testing.T) {
		// Given
		suite.SetupTest()
		route := schema.Route{
			Path:   "/v1/test",
			Method: "GET",
		}
		routeProperty := schema.RouteProperty{}
		suite.routeManager.AddRoute(route, routeProperty)

		// When
		returnedRouteProperty, err := suite.routeManager.GetRouteProperty(route)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, routeProperty, returnedRouteProperty)
	})

	suite.T().Run("cannot get route property", func(t *testing.T) {
		// Given
		suite.SetupTest()
		route := schema.Route{
			Path:   "/v1/test",
			Method: "GET",
		}

		// When
		_, err := suite.routeManager.GetRouteProperty(route)

		// Then
		assert.Error(t, err)

	})

}

func (suite *DolusRouteManagerTestSuite) TestGetRouteProperties() {

	suite.T().Run("can get route properties", func(t *testing.T) {
		// Given
		suite.SetupTest()
		routes := []schema.Route{
			{Path: "/v1/test", Method: "GET"},
			{Path: "/v1/test", Method: "POST"},
			{Path: "/v2/test", Method: "GET"},
		}
		routeProperty := schema.RouteProperty{
			RequestSchema: dstruct.NewGeneratedStruct(struct {
				Age int
			}{}),
		}

		for _, route := range routes {
			suite.routeManager.AddRoute(route, routeProperty)
		}

		// When
		routeProperties := suite.routeManager.GetRouteProperties()

		// Then
		assert.Equal(t, len(routes), len(routeProperties))
		for _, route := range routes {
			rp, ok := routeProperties[route]
			assert.True(t, ok)
			assert.Equal(t, routeProperty, rp)
		}

	})
}

func TestRouteManagerImplTestSuite(t *testing.T) {
	suite.Run(t, new(DolusRouteManagerTestSuite))
}
