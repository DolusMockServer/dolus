package api

import "github.com/DolusMockServer/dolus/pkg/schema"

type DolusApi interface {
	ServerInterface
	AddRoute(pathMethod schema.Route) error
}
