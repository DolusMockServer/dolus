package engine

import (
	"net/http"

	"github.com/MartinSimango/dolus/core"
	"github.com/MartinSimango/dolus/expectation"
)

type ExpectationEngine interface {
	AddResponseSchemaForPathMethod(responseSchema *core.ResponseSchema) error
	AddExpectationsFromFiles(files ...string)
	AddExpectation(pathMethod expectation.PathMethod, expectation expectation.Expectation) error
	GetExpectations() map[expectation.PathMethod][]expectation.Expectation
	GetExpectationForPathMethod(pathMethod expectation.PathMethod) []expectation.Expectation
	GetResponseForRequest(path, method string, request *http.Request) (*expectation.Response, error)
	Load() error
}

// GetEchoResponse(path, method string, ctx echo.Context) error
