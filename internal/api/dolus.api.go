package api

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type DolusApi interface {
	ServerInterface
	AddRoute(pathMethod expectation.Route) error
}
