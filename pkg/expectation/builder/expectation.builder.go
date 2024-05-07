package builder

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
)

type Output struct {
	Expectations []expectation.Expectation
	RouteManager engine.RouteManager
}

type ExpectationBuilder interface {
	BuildExpectations() (*Output, error)
}
