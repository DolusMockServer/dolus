package engine

import (
	"reflect"
	"slices"
	"testing"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type DolusExpectationEngineTestSuite struct {
	suite.Suite
	expectationEngine ExpectationEngine
	routeManager      *RouteManagerMock
}

func (suite *DolusExpectationEngineTestSuite) SetupTest() {
	suite.routeManager = NewRouteManagerMock(suite.T())
	suite.expectationEngine = NewDolusExpectationEngine(*generator.NewGenerationConfig(), suite.routeManager, nil)
}

func (suite *DolusExpectationEngineTestSuite) TestAddExpectation() {
	suite.T().Run("Can add an expectation", func(t *testing.T) {
		// Given
		exp := expectation.Expectation{}

		// When
		err := suite.expectationEngine.AddExpectation(exp)

		// Then
		returnedExpectations := suite.expectationEngine.GetExpectations(
			toPointer(exp.ExpectationType),
			toPointer(exp.Request.Path),
			toPointer(exp.Request.Method),
		)
		assert.NoError(t, err)
		assert.True(
			t,
			slices.ContainsFunc(returnedExpectations, func(e expectation.Expectation) bool {
				return reflect.DeepEqual(e, exp)
			}),
		)
	})
}

func (suite *DolusExpectationEngineTestSuite) TestAddExpectations() {
}

func (suite *DolusExpectationEngineTestSuite) TestGetExpectations() {
}

func (suite *DolusExpectationEngineTestSuite) TestGetResponseForRequest() {
}

func TestDolusExpectationEngineTestSuite(t *testing.T) {
	suite.Run(t, new(DolusExpectationEngineTestSuite))
}

func toPointer[T any](v T) *T {
	return &v
}
