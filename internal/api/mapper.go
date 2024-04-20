package api

import "github.com/DolusMockServer/dolus/pkg/expectation"

type Mapper interface {
	MapCueExpectations(
		expectations []expectation.Expectation,
	) ([]Expectation, error)
	MapCueExpectation(expectation expectation.Expectation) (*Expectation, error)
}
