package builder

import "github.com/DolusMockServer/dolus/pkg/expectation"

type ExpectationBuilder interface {
	BuildExpectations() ([]expectation.DolusExpectation, error)
}
