package builder

import (
	"github.com/DolusMockServer/dolus/pkg/expectation/models"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type Output struct {
	Expectations    []models.Expectation
	RouteProperties schema.RouteProperties
}

type ExpectationBuilder interface {
	BuildExpectations() (*Output, error)
}
