package api

// integration tests for the dolus api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
)

/*

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
	expectationEngine *engine.ExpectationEngineMock
	mapper            *MapperMock
	api               DolusApi
}

func (suite *DolusApiImplTestSuite) SetupTest() {
	suite.expectationEngine = engine.NewExpectationEngineMock(suite.T())
	suite.mapper = NewMapperMock(suite.T())
	suite.api = NewDolusApi(suite.expectationEngine, suite.mapper)
}

func (suite *DolusApiImplTestSuite) TestGetV1DolusExpectations() {
	suite.T().Run("should return 200 OK with expectations", func(t *testing.T) {
		// Given
		suite.SetupTest()
		expectations := expectation.Expectations{
			Expectations: []expectation.Expectation{},
		}
		suite.expectationEngine.EXPECT().GetExpectations(nil, nil, nil).Return(expectations.Expectations)
		suite.mapper.EXPECT().
			MapToApiExpectations(expectations.Expectations).
			Return([]Expectation{}, nil)
		req := httptest.NewRequest(http.MethodGet, "/v1/dolus/expectations", nil)
		rec := httptest.NewRecorder()

		// When
		err := suite.api.GetExpectations(echo.New().NewContext(req, rec), GetExpectationsParams{})

		// Then
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	suite.T().Run("should return 500 if mapper fails", func(t *testing.T) {
		// Given
		suite.SetupTest()
		expectations := expectation.Expectations{}
		suite.expectationEngine.EXPECT().GetExpectations(nil, nil, nil).Return(expectations.Expectations)
		suite.mapper.EXPECT().MapToApiExpectations(expectations.Expectations).Return(nil,
			fmt.Errorf("error"))
		req := httptest.NewRequest(http.MethodGet, "/v1/dolus/expectations", nil)
		rec := httptest.NewRecorder()

		// When
		err := suite.api.GetExpectations(echo.New().NewContext(req, rec), GetExpectationsParams{})

		// Then
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func (suite *DolusApiImplTestSuite) TestCreateDolusExpectations() {
	suite.T().Run("should return 201 created if request is successful", func(t *testing.T) {
		// Given
		suite.SetupTest()

		var request Expectation
		var mappedRequest *expectation.Expectation = &expectation.Expectation{
			Priority: 1,
		}
		suite.mapper.EXPECT().MapToExpectation(request).Return(mappedRequest, nil)
		suite.expectationEngine.EXPECT().AddExpectation(mappedRequest, true).Return(nil)

		requestBody, err := json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/dolus/expectations",
			bytes.NewBuffer(requestBody),
		)
		rec := httptest.NewRecorder()

		// When
		resultErr := suite.api.CreateExpectation(echo.New().NewContext(req, rec))
		var result Expectation
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		if err != nil {
			t.Fatal(err)
		}

		// Then
		assert.NoError(t, resultErr)

		assert.Equal(t, http.StatusCreated, rec.Code)

		assert.Equal(t, 1, result.Priority)
	})

	suite.T().Run("should return 400 if malformed function", func(t *testing.T) {
	})

	suite.T().Run("should return 500", func(t *testing.T) {
	})
}

func TestDolusApiTestSuite(t *testing.T) {
	suite.Run(t, new(DolusApiImplTestSuite))
}
