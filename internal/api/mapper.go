package api

import (
	"github.com/DolusMockServer/dolus/pkg/expectation/cue"
)

type Mapper interface {
	MapCueExpectations(
		expectations []cue.Expectation,
	) ([]Expectation, error)
	MapCueExpectation(expectation cue.Expectation) (*Expectation, error)
}
