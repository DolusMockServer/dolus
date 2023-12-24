package engine

import (
	"net/http"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type ExpectationEngine interface {
	AddExpectation(expectation expectation.DolusExpectation, validateExpectationSchema bool) error
	AddResponseSchemaForPathMethodStatus(
		pathMethodStatus expectation.PathMethodStatus,
		schema dstruct.DynamicStructModifier,
	) error
	GetExpectations() map[string][]expectation.DolusExpectation
	GetExpectationForPathMethod(pathMethod expectation.PathMethod) []expectation.DolusExpectation
	GetResponseForRequest(
		path, method string,
		request *http.Request,
	) (*expectation.DolusResponse, error)
	GetRawCueExpectations() dolus.Expectations
	GetExpectationRoutes() []string
}
