package api

import (
	"github.com/DolusMockServer/dolus/internal/server"
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type DolusApi interface {
	server.ServerInterface
	AddRoute(pathMethod expectation.Route) error
}
