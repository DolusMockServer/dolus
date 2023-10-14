package dolus

import (
	"fmt"

	"github.com/MartinSimango/dolus/engine"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
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
	website = "https://github.com/MartinSimango/dolus"
)

func printBanner() {
	versionColor := color.New(color.FgGreen).SprintFunc()("v", Version)
	websiteColor := color.New(color.FgBlue).SprintFunc()(website)
	fmt.Printf(banner, versionColor, websiteColor)
}

type Dolus struct {
	OpenAPIspec        string
	HideBanner         bool
	HidePort           bool
	EchoServer         *echo.Echo
	GenerationConfig   generator.GenerationConfig
	expectationEngine  engine.ExpectationEngine
	expectationFiles   []string
	openAPISpecLoader  expectation.Loader[expectation.OpenAPISpecLoadType]
	cueLoader          expectation.Loader[expectation.CueExpectationLoadType]
	expectationBuilder expectation.ExpectationBuilder
	fieldGenerator     *generator.Generator
}

func New() *Dolus {
	generationConfig := generator.NewGenerationConfig()

	return &Dolus{
		HideBanner:       false,
		HidePort:         false,
		OpenAPIspec:      "openapi.yaml",
		GenerationConfig: *generationConfig,
	}

}

func (d *Dolus) initHttpServer() {
	d.EchoServer = echo.New()
	d.EchoServer.HideBanner = true
	d.EchoServer.HidePort = d.HidePort
	d.openAPISpecLoader = expectation.NewOpenOPISpecLoader(d.OpenAPIspec)
	d.cueLoader = expectation.NewCueExpectationLoader(d.expectationFiles)
	d.fieldGenerator = generator.NewGenerator(&d.GenerationConfig)
	d.expectationBuilder = expectation.NewExpectationBuilderImpl(*d.fieldGenerator)

}

func (d *Dolus) addRoutes(method, path string) {
	d.EchoServer.Router().Add(method, path, func(ctx echo.Context) error {
		response, err := d.expectationEngine.GetResponseForRequest(path, method, ctx.Request())
		if err != nil {
			return ctx.JSON(500, GeneralError{
				Path:     ctx.Request().URL.Path,
				Method:   method,
				ErrorMsg: err.Error(),
			})
		}

		response.Body.Generate()
		response.Body.Update()
		return ctx.JSON(response.Status, response.Body.Instance())
	})
}

func (d *Dolus) loadOpenAPISpecExpectations() error {

	expectations, err := d.expectationBuilder.BuildExpectationsFromOpenApiSpecLoader(d.openAPISpecLoader)
	if err != nil {
		return err
	}

	for _, e := range expectations {

		method := e.Request.Method
		path := e.Request.Path
		d.addRoutes(method, path)
		d.expectationEngine.AddResponseSchemaForPathMethodStatus(expectation.PathMethodStatusExpectation(e),
			e.Response.Body)

		if err := d.expectationEngine.AddExpectation(e, false); err != nil {
			fmt.Printf("Error adding expectation:\n%s\n", err)
		}

	}

	return nil
}

func (d *Dolus) loadCueExpectations() error {

	expectations, err := d.expectationBuilder.BuildExpectationsFromCueLoader(d.cueLoader)
	if err != nil {
		return err
	}
	for _, e := range expectations {
		d.expectationEngine.AddExpectation(e, true)

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
	// d.expectationEngine.Load()
	if err := d.loadExpectations(); err != nil {
		return err
	}

	return d.EchoServer.Start(address)
}

func (d *Dolus) Start(address string) error {
	if !d.HideBanner {
		printBanner()
	}

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

}
