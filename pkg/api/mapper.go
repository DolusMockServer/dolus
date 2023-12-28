package api

import (
	"github.com/DolusMockServer/dolus/internal/server"
	"github.com/DolusMockServer/dolus/pkg/expectation/cue"
)

type Mapper interface {
	MapCueExpectations(
		expectations []cue.Expectation,
	) ([]server.Expectation, error)
}
