package dolus

import (
	"fmt"
	"net/http"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/DolusMockServer/dolus/internal/server"
	"github.com/DolusMockServer/dolus/pkg/api"
	"github.com/DolusMockServer/dolus/pkg/engine"
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/builder"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/task"
)

const (
	banner = `  
    ____            __                
   / __ \  ____    / / __  __   _____
  / / / / / __ \  / / / / / /  / ___/
 / /_/ / / /_/ / / / / /_/ /  (__  ) 
/_____/  \____/ /_/  \____/  /____/ %s
Go framework for creating customizable and extendable mock servers
%s

--------------------------------------------------------------------

`
	Version = "0.0.1"
	website = "https://github.com/DolusMockServer/dolus"
)

func printBanner() {
	logger.Log.SetTextFormatter()
	versionColor := color.New(color.FgGreen).SprintFunc()("v", Version)
	websiteColor := color.New(color.FgBlue).SprintFunc()(website)
	logger.Log.Infof(banner, versionColor, websiteColor)
	logger.Log.SetDefaultFormatter()
}

type Dolus struct {
	OpenAPIspec               string
	HideBanner                bool
	HidePort                  bool
	EchoServer                *echo.Echo
	GenerationConfig          generator.GenerationConfig
	expectationEngine         engine.ExpectationEngine
	expectationFiles          []string
	cueExpectationBuilder     builder.ExpectationBuilder
	openApiExpectationBuilder builder.ExpectationBuilder
	fieldGenerator            *generator.Generator
	dolusApi                  api.DolusApi
}

func New() *Dolus {
	logger.Log = logger.NewLogger("logfile.log")
	generationConfig := generator.NewGenerationConfig()
	return &Dolus{
		HideBanner:       false,
		HidePort:         false,
		OpenAPIspec:      "openapi.yaml",
		GenerationConfig: *generationConfig,
	}
}

func (d *Dolus) initMiddleware() error {
	d.EchoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (d *Dolus) initHttpServer() {
	d.EchoServer = echo.New()
	d.EchoServer.HideBanner = true
	d.EchoServer.HidePort = d.HidePort
	d.fieldGenerator = generator.NewGenerator(&d.GenerationConfig)
	d.cueExpectationBuilder = builder.NewCueExpectationBuilder(
		d.expectationFiles,
		*d.fieldGenerator,
	)
	d.openApiExpectationBuilder = builder.NewOpenApiExpectationBuilder(
		d.OpenAPIspec,
		*d.fieldGenerator,
	)
	d.dolusApi = api.NewDolusApi(d.expectationEngine, api.NewMapper())

	d.initMiddleware()

	server.RegisterHandlers(d.EchoServer, d.dolusApi)
}

// TODO: move this to somewhere more appropriate
type GeneralError struct {
	Path     string
	Method   string
	ErrorMsg string
}

func (d *Dolus) addRoutes(route expectation.Route) {
	if err := d.dolusApi.AddRoute(route); err != nil {
		fmt.Printf("error adding route: %s\n", err.Error())
	}
	d.EchoServer.Router().Add(route.Method, route.Path, func(ctx echo.Context) error {
		logger.Log.Infof(
			"Received request for path %s and method %s",
			ctx.Request().URL.Path,
			route.Method,
		)
		response, err := d.expectationEngine.GetResponseForRequest(route.Path, ctx.Request())
		if err != nil {
			return ctx.JSON(500, GeneralError{
				Path:     ctx.Request().URL.Path,
				Method:   route.Method,
				ErrorMsg: err.Error(),
			})
		}

		response.Body.GenerateAndUpdate()
		return ctx.JSON(response.Status, response.Body.Instance())
	})
}

func (d *Dolus) loadOpenAPISpecExpectations() error {
	expectations, err := d.openApiExpectationBuilder.BuildExpectations()
	if err != nil {
		return err
	}
	for _, e := range expectations {
		d.addRoutes(e.Request.Route)
		d.expectationEngine.AddResponseSchemaForRoute(
			e.Request.Route,
			e.Response.Body,
		)

		if err := d.expectationEngine.AddExpectation(e, false); err != nil {
			fmt.Printf("Error adding expectation:\n%s\n", err)
		}

	}

	return nil
}

func (d *Dolus) loadCueExpectations() error {
	expectations, err := d.cueExpectationBuilder.BuildExpectations()
	if err != nil {
		return err
	}
	for _, e := range expectations {
		if err := d.expectationEngine.AddExpectation(e, true); err != nil {
			fmt.Printf("Error adding expectation:\n%s\n", err)
		}
	}

	return nil
}

func (d *Dolus) loadExpectations() error {
	if err := d.loadOpenAPISpecExpectations(); err != nil {
		return err
	}

	if err := d.loadCueExpectations(); err != nil {
		return err
	}
	return nil
}

func (d *Dolus) startHttpServer(address string) error {
	d.initHttpServer()
	go task.RegisterDolusTasks()
	if err := d.loadExpectations(); err != nil {
		return err
	}
	d.EchoServer.Logger.SetOutput(logger.Log.Out)
	return d.EchoServer.Start(address)
}

func (d *Dolus) Start(address string) error {
	if !d.HideBanner {
		printBanner()
	}
	logger.Log.SetLevel(logrus.InfoLevel)

	if d.expectationEngine == nil {
		generationConfig := d.GenerationConfig
		generationConfig.SetNonRequiredFields(true)
		d.expectationEngine = engine.NewDolusExpectationEngine(generationConfig)
	}

	return d.startHttpServer(address)
}

func (d *Dolus) AddExpectations(files ...string) {
	d.expectationFiles = append(d.expectationFiles, files...)
}

func (d *Dolus) AddCueExpectationsFromFolder() {
	// TODO: implement
}
