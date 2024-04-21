package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

/*
  // Fetch expectations
    // (GET /v1/dolus/expectations)
    GetV1DolusExpectations(ctx echo.Context) error

    // (POST /v1/dolus/expectations)
    PostV1DolusExpectations(ctx echo.Context) error
    // Your GET endpoint
    // (GET /v1/dolus/logs)
    GetV1DolusLogs(ctx echo.Context, params GetV1DolusLogsParams) error
    // Your GET endpoint
    // (GET /v1/dolus/logs/ws)
    GetV1DolusLogsWs(ctx echo.Context, params GetV1DolusLogsWsParams) error
    // Your GET endpoint
    // (GET /v1/dolus/routes)
    GetV1DolusRoutes(ctx echo.Context) error

*/

type DolusApiImplTestSuite struct {
	suite.Suite
	mockEngine *engine.ExpectationEngineMock
	mockMapper *MapperMock
	api        DolusApi
}

func (suite *DolusApiImplTestSuite) SetupTest() {
	fmt.Println("RESET!")
	suite.mockEngine = engine.NewExpectationEngineMock(suite.T())
	suite.mockMapper = NewMapperMock(suite.T())
	suite.api = NewDolusApi(suite.mockEngine, suite.mockMapper)
}

func (suite *DolusApiImplTestSuite) TestGetV1DolusExpectations() {

	suite.T().Run("should return 200 OK with expectations", func(t *testing.T) {
		expectations := expectation.Expectations{}
		suite.mockEngine.EXPECT().GetCueExpectations().Return(expectations)
		suite.mockMapper.EXPECT().MapCueExpectations(expectations.Expectations).Return([]Expectation{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/v1/dolus/expectations", nil)
		rec := httptest.NewRecorder()

		err := suite.api.GetV1DolusExpectations(echo.New().NewContext(req, rec))
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusOK)
	})

	suite.T().Run("should return 500", func(t *testing.T) {
		suite.SetupTest()
		expectations := expectation.Expectations{}
		suite.mockEngine.EXPECT().GetCueExpectations().Return(expectations)
		suite.mockMapper.EXPECT().MapCueExpectations(expectations.Expectations).Return(nil,
			fmt.Errorf("error"))

		req := httptest.NewRequest(http.MethodGet, "/v1/dolus/expectations", nil)
		rec := httptest.NewRecorder()

		err := suite.api.GetV1DolusExpectations(echo.New().NewContext(req, rec))
		assert.NoError(t, err)
		assert.Equal(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestDolusApiTestSuite(t *testing.T) {
	suite.Run(t, new(DolusApiImplTestSuite))
}
