package api

import (
	"github.com/DolusMockServer/dolus/pkg/expectation/models"
)

type Mapper interface {
	MapCueExpectations(
		expectations []models.Expectation,
	) ([]Expectation, error)
	MapCueExpectation(expectation models.Expectation) (*Expectation, error)
}
