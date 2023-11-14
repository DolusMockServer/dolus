package engine

import (
	"net/http"

	"github.com/MartinSimango/dolus-expectations/pkg/dolus"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dstruct"
)

type ExpectationEngine interface {
	AddExpectation(expectation expectation.Expectation, validateExpectationSchema bool) error
	AddResponseSchemaForPathMethodStatus(pathMethodStatus expectation.PathMethodStatus, schema dstruct.DynamicStructModifier) error
	GetExpectations() map[expectation.PathMethod][]expectation.Expectation
	GetExpectationForPathMethod(pathMethod expectation.PathMethod) []expectation.Expectation
	GetResponseForRequest(path, method string, request *http.Request) (*expectation.Response, error)
	GetRawCueExpectations() dolus.Expectations
	GetExpectationRoutes() []string
}
