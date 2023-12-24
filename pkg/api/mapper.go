package api

import (
	cueDolus "github.com/DolusMockServer/dolus-expectations/pkg/dolus"
	"github.com/DolusMockServer/dolus/internal/server"
)

type Mapper interface {
	MapCueExpectations(
		expectations []cueDolus.Expectation,
	) ([]server.Expectation, error)
}
