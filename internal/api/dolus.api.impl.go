package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type DolusApiImpl struct {
	ExpectationEngine engine.ExpectationEngine
	Mapper            Mapper
	routes            map[schema.Route]bool
}

var _ DolusApi = &DolusApiImpl{}

func NewDolusApi(expectationEngine engine.ExpectationEngine,
	mapper Mapper,
) *DolusApiImpl {
	return &DolusApiImpl{
		ExpectationEngine: expectationEngine,
		Mapper:            mapper,
		routes:            make(map[schema.Route]bool),
	}
}

func (d *DolusApiImpl) AddRoute(route schema.Route) error {
	if d.routes[route] {
		return fmt.Errorf(
			"route %s with operation %s already exists",
			route.Path,
			route.Method,
		)
	}
	d.routes[route] = true
	return nil
}

// GetV1DolusExpectations implements server.ServerInterface.
func (d *DolusApiImpl) GetV1DolusExpectations(ctx echo.Context) error {
	apiExpectations, err := d.Mapper.MapCueExpectations(
		d.ExpectationEngine.
			GetCueExpectations().
			Expectations)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return ctx.JSON(http.StatusOK,
		apiExpectations)
}

// GetV1DolusRoutes implements server.ServerInterface.
func (d *DolusApiImpl) GetV1DolusRoutes(ctx echo.Context) error {
	var serverRoutes []Route
	for r := range d.routes {
		serverRoutes = append(serverRoutes, Route{
			Path:      r.Path,
			Operation: r.Method,
		})
	}

	return ctx.JSON(200, serverRoutes)
}

// PostV1DolusExpectations implements server.ServerInterface.
func (d *DolusApiImpl) PostV1DolusExpectations(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
}

// GetV1DolusLogs implements DolusApi.
func (*DolusApiImpl) GetV1DolusLogs(ctx echo.Context, params GetV1DolusLogsParams) error {
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
func (*DolusApiImpl) GetV1DolusLogsWs(
	ctx echo.Context,
	params GetV1DolusLogsWsParams,
) error {
	return ctx.JSON(http.StatusNotImplemented, "Not Implemented")
	// conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	// if err != nil {
	// 	return err
	// }
	// logger.Log.RegisterWebSocketClient(conn, 1000)
	// return nil
}
