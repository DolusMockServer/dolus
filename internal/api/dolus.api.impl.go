package api

import (
	"fmt"
	"net/http"

	"github.com/MartinSimango/dstruct"
	"github.com/labstack/echo/v4"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type DolusApiImpl struct {
	ExpectationEngine engine.ExpectationEngine
	Mapper            Mapper
}

var _ DolusApi = &DolusApiImpl{}

func NewDolusApi(expectationEngine engine.ExpectationEngine,
	mapper Mapper,
) *DolusApiImpl {
	return &DolusApiImpl{
		ExpectationEngine: expectationEngine,
		Mapper:            mapper,
	}
}

// GetExpectations implements server.ServerInterface.
func (d *DolusApiImpl) GetExpectations(ctx echo.Context, params GetExpectationsParams) error {

	var expectationType *expectation.ExpectationType
	var method *string
	var ok bool

	if params.ExpectationType != nil {
		expectationType = new(expectation.ExpectationType)
		*expectationType, ok = any(*params.ExpectationType).(expectation.ExpectationType)
		if !ok {
			return ctx.JSON(http.StatusBadRequest, BadRequest{
				Message: fmt.Sprintf("invalid expectation type: %s", *params.ExpectationType),
			})
		}
	}
	if params.Method != nil {
		method = new(string)
		*method, ok = any(*params.Method).(string)
		if !ok {
			return ctx.JSON(http.StatusBadRequest, BadRequest{
				Message: fmt.Sprintf("invalid method: %s", *params.Method),
			})
		}
	}

	apiExpectations, err := d.Mapper.MapToApiExpectations(
		d.ExpectationEngine.GetExpectations(expectationType, params.Path, method))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK,
		apiExpectations)
}

// GetV1DolusRoutes implements server.ServerInterface.
func (d *DolusApiImpl) GetRoutes(ctx echo.Context) error {
	var serverRoutes []Route
	for _, r := range d.ExpectationEngine.GetRoutes() {
		serverRoutes = append(serverRoutes, Route{
			Path:      r.Path,
			Operation: r.Method,
		})
	}

	return ctx.JSON(200, serverRoutes)
}

// PostV1DolusExpectations implements server.ServerInterface.
func (d *DolusApiImpl) CreateExpectation(ctx echo.Context) error {

	defer ctx.Request().Body.Close()
	var apiExpectation Expectation

	if err := ctx.Bind(&apiExpectation); err != nil {
		echoError := err.(*echo.HTTPError)
		return ctx.JSON(echoError.Code, BadRequest{
			Message: echoError.Internal.Error(),
		})
	}

	expct, err := d.Mapper.MapToExpectation(apiExpectation)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequest{
			Message: err.Error(),
		})
	}

	expct.ExpectationType = expectation.Custom

	if expct.Response.Body != nil {
		expct.Response.GeneratedBody = dstruct.NewGeneratedStruct(schema.SchemaFromAny(expct.Response.Body))
	}

	// respsonse :=
	// oeb.fieldGenerator.GenerationConfig.SetValueGenerationType(generator.UseDefaults)

	if err := d.ExpectationEngine.AddExpectation(*expct); err != nil {
		// TODO: depending on the error, return a different status code
		return ctx.JSON(http.StatusInternalServerError, InternalServerError{
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusCreated, expct)
}

// GetV1DolusLogs implements DolusApi.
func (*DolusApiImpl) GetLogs(ctx echo.Context, params GetLogsParams) error {
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
