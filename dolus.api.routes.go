package dolus

import (
	"net/http"

	"github.com/MartinSimango/dolus/engine"
	"github.com/MartinSimango/dolus/server"
	echo "github.com/labstack/echo/v4"
)

type DolusApiRoutes struct {
	ExpectationEngine engine.ExpectationEngine
	DolusApiFactory   DolusApiFactory
}

var _ server.ServerInterface = &DolusApiRoutes{}

func NewDolusApiRoutes(expectationEngine engine.ExpectationEngine,
	dolusApiFactory DolusApiFactory,
) *DolusApiRoutes {
	return &DolusApiRoutes{
		ExpectationEngine: expectationEngine,
		DolusApiFactory:   dolusApiFactory,
	}
}

// GetV1DolusExpectations implements server.ServerInterface.
func (d *DolusApiRoutes) GetV1DolusExpectations(ctx echo.Context) error {
	apiExpectations, err := d.DolusApiFactory.RawCueExpectationsToApiExpectations(d.ExpectationEngine.GetRawCueExpectations().Expectations)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return ctx.JSON(http.StatusOK,
		apiExpectations)
}

// GetV1DolusRoutes implements server.ServerInterface.
func (*DolusApiRoutes) GetV1DolusRoutes(ctx echo.Context) error {
	panic("unimplemented")
}

// PostV1DolusExpectations implements server.ServerInterface.
func (*DolusApiRoutes) PostV1DolusExpectations(ctx echo.Context) error {
	panic("unimplemented")
}
