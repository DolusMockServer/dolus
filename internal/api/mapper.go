package api

import "github.com/DolusMockServer/dolus/pkg/expectation"

type Mapper interface {
	MapToApiExpectations(
		expectations []expectation.Expectation,
	) ([]Expectation, error)
	MapToApiExpectation(expectation expectation.Expectation) (*Expectation, error)
	MapToExpectation(expectation Expectation) (*expectation.Expectation, error)
}
