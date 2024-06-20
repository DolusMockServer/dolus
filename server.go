package dolus

import (
	"fmt"
	"time"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/DolusMockServer/dolus/internal/api"
	"github.com/DolusMockServer/dolus/pkg/expectation/builder"
	"github.com/DolusMockServer/dolus/pkg/expectation/engine"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
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

type Server struct {
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
	schemaMapper              schema.Mapper
}

func New() *Server {
	logger.Log = logger.NewLogger("logfile.log")

	generationConfig := generator.NewGenerationConfig()
	return &Server{
		HideBanner:       false,
		HidePort:         false,
		OpenAPIspec:      "openapi.yaml",
		GenerationConfig: *generationConfig,
		schemaMapper:     schema.NewMapper(),
	}
}

func printBanner() {
	logger.Log.SetTextFormatter()
	versionColor := color.New(color.FgGreen).SprintFunc()("v", Version)
	websiteColor := color.New(color.FgBlue).SprintFunc()(website)
	logger.Log.Infof(banner, versionColor, websiteColor)
	logger.Log.SetDefaultFormatter()
}

func (d *Server) initMiddleware() error {
	d.EchoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	return nil
}

func (d *Server) initHttpServer() {
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

	api.RegisterHandlers(d.EchoServer, d.dolusApi)
}

func (d *Server) addHandlersForRoutes(routes []schema.Route) {
	for _, route := range routes {
		d.EchoServer.Router().Add(route.Method, route.Path, func(ctx echo.Context) error {
			logger.Log.Infof(
				"Received request for path %s and method %s",
				ctx.Request().URL.RequestURI(),
				route.Method,
			)

			response, err := d.expectationEngine.GetResponseForRequest(
				ctx.Request(),
				d.schemaMapper.MapToRequestParameters(ctx),
				route.Path,
			)
			if err != nil {
				return ctx.JSON(500, api.GeneralError{
					Path:     ctx.Request().URL.Path,
					Method:   route.Method,
					ErrorMsg: err.Error(),
				})
			}

			response.GeneratedBody.GenerateAndUpdate()
			return ctx.JSON(response.Status, response.GeneratedBody.Instance())
		})
	}
}

func (d *Server) loadOpenAPISpecExpectations() error {
	output, err := d.openApiExpectationBuilder.BuildExpectations()
	if err != nil {
		return err
	}

	d.expectationEngine.SetRouteManager(output.RouteManager)
	d.addHandlersForRoutes(output.RouteManager.GetRoutes())

	for _, e := range output.Expectations {
		if err := d.expectationEngine.AddExpectation(e, false); err != nil {
			fmt.Printf("Error adding expectation:\n%s\n", err)
		}
	}

	return nil
}

func (d *Server) loadCueExpectations() error {
	output, err := d.cueExpectationBuilder.BuildExpectations()
	if err != nil {
		return err
	}
	for _, e := range output.Expectations {
		if err := d.expectationEngine.AddExpectation(e, true); err != nil {
			fmt.Printf("Error adding expectation:\n%s\n", err)
		}
	}

	return nil
}

func (d *Server) loadExpectations() error {
	s := time.Now()
	if err := d.loadOpenAPISpecExpectations(); err != nil {
		return err
	}
	fmt.Println("Time to load openapi expectations: ", time.Since(s))
	s2 := time.Now()
	if err := d.loadCueExpectations(); err != nil {
		return err
	}
	fmt.Println("Time to load cue expectations: ", time.Since(s2))
	return nil
}

func (d *Server) startHttpServer(address string) error {
	d.initHttpServer()
	go task.RegisterDolusTasks()
	now := time.Now()

	if err := d.loadExpectations(); err != nil {
		return err
	}
	logger.Log.Infof("Server started in %s", time.Since(now))

	d.EchoServer.Logger.SetOutput(logger.Log.Out)
	return d.EchoServer.Start(address)
}

func (d *Server) Start(address string) error {
	if !d.HideBanner {
		printBanner()
	}
	logger.Log.SetLevel(logrus.DebugLevel)

	if d.expectationEngine == nil {
		generationConfig := d.GenerationConfig
		generationConfig.SetNonRequiredFields(true)
		d.expectationEngine = engine.NewDolusExpectationEngine(generationConfig)
	}

	return d.startHttpServer(address)
}

func (d *Server) AddExpectations(files ...string) {
	d.expectationFiles = append(d.expectationFiles, files...)
}

func (d *Server) AddCueExpectationsFromFolder() {
}
