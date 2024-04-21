package api

import (
	"testing"
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

func TestGetV1DolusExpectations(t *testing.T) {
	// TODO mock mapper
	// mockEngine := new(engine.MockEngine)
	mockMapper := api.NewMockMapper(t)
	mockMapper

	// Create new DolusApi instance
	// api := NewDolusApi(mockEngine, mockMapper)

	// Define a route
	// route := schema.Route{
	// 	Path:   "/test",
	// 	Method: "GET",
	// }

	// // Test AddRoute method
	// err := api.AddRoute(route)
	// assert.NoError(t, err)

	// // Test adding the same route again
	// err = api.AddRoute(route)
	// assert.Error(t, err)

	// // Test GetV1DolusExpectations method
	// ctx := echo.New().NewContext(nil, nil)
	// err = api.GetV1DolusExpectations(ctx)
	// assert.NoError(t, err)
}
