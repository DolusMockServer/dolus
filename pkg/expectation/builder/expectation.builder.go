package builder

import (
	"strings"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

type ExpectationBuilder interface {
	BuildExpectations() ([]expectation.DolusExpectation, error)
}

func pathFromOpenApiPath(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}
