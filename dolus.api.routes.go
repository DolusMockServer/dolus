package dolus

import (
	"fmt"
	"net/http"

	"github.com/MartinSimango/dolus/engine"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dolus/logger"
	"github.com/MartinSimango/dolus/server"
	echo "github.com/labstack/echo/v4"
)

type DolusApi interface {
	server.ServerInterface
	AddRoute(pathMethod expectation.PathMethod) error
}

type DolusApiRoutes struct {
	ExpectationEngine engine.ExpectationEngine
	DolusApiFactory   DolusApiFactory
	routes            map[expectation.PathMethod]bool
}

var _ DolusApi = &DolusApiRoutes{}

func NewDolusApiRoutes(expectationEngine engine.ExpectationEngine,
	dolusApiFactory DolusApiFactory,

) *DolusApiRoutes {
	return &DolusApiRoutes{
		ExpectationEngine: expectationEngine,
		DolusApiFactory:   dolusApiFactory,
		routes:            make(map[expectation.PathMethod]bool),
	}
}

func (d *DolusApiRoutes) AddRoute(pathMethod expectation.PathMethod) error {
	if d.routes[pathMethod] {
		return fmt.Errorf("route %s with operation %s already exists", pathMethod.Path, pathMethod.Method)
	}
	d.routes[pathMethod] = true
	return nil
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
func (d *DolusApiRoutes) GetV1DolusRoutes(ctx echo.Context) error {
	var serverRoutes []server.Route
	for k := range d.routes {
		serverRoutes = append(serverRoutes, server.Route{
			Path:      k.Path,
			Operation: k.Method,
		})

	}

	return ctx.JSON(200, serverRoutes)

}

// PostV1DolusExpectations implements server.ServerInterface.
func (*DolusApiRoutes) PostV1DolusExpectations(ctx echo.Context) error {
	panic("unimplemented")
}

// GetV1DolusLogs implements DolusApi.
func (*DolusApiRoutes) GetV1DolusLogs(ctx echo.Context, params server.GetV1DolusLogsParams) error {
	lines := 1000
	if params.Lines != nil {
		lines = *params.Lines
	}
	if logs, err := logger.Log.GetLogStream(lines); err != nil {
		return ctx.JSON(http.StatusInternalServerError, struct {
			Message string
		}{Message: err.Error()})
	} else {
		return ctx.String(http.StatusOK, logs)
	}
}

// GetV1DolusLogsWs implements DolusApi.
func (*DolusApiRoutes) GetV1DolusLogsWs(ctx echo.Context, params server.GetV1DolusLogsWsParams) error {
	panic("Implementation doesn't work")
	conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	if err != nil {
		return err
	}
	logger.Log.RegisterWebSocketClient(conn, 1000)
	return nil
}
