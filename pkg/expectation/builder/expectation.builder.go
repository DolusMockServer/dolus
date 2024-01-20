package builder

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type Output struct {
	Expectations    []expectation.Expectation
	RouteProperties schema.RouteProperties
}

type ExpectationBuilder interface {
	BuildExpectations() (*Output, error)
}
