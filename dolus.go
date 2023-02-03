package dolus

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MartinSimango/dolus/core"
	"github.com/MartinSimango/dolus/engine"
	"github.com/MartinSimango/dolus/expectation"
	"github.com/MartinSimango/dolus/generator"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
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
	OpenAPIspec       string
	HideBanner        bool
	HidePort          bool
	EchoServer        *echo.Echo
	GenerationConfig  generator.GenerationConfig
	expectationEngine engine.ExpectationEngine
	expectationFiles  []string
}

func New() *Dolus {
	return &Dolus{
		HideBanner:       false,
		HidePort:         false,
		OpenAPIspec:      "openapi.yaml",
		GenerationConfig: *generator.NewGenerationConfig(),
	}

}

func (d *Dolus) initHttpServer() {
	d.EchoServer = echo.New()
	d.EchoServer.HideBanner = true
	d.EchoServer.HidePort = d.HidePort
}

func getRealPath(path string) string {
	p := strings.ReplaceAll(path, "{", ":")
	return strings.ReplaceAll(p, "}", "")
}

func (d *Dolus) startHttpServer(address string) error {
	start := time.Now()
	d.initHttpServer()

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(d.OpenAPIspec)
	if err != nil {
		return err
	}

	if err := doc.Validate(ctx); err != nil {
		return err
	}
	for path := range doc.Paths {
		for method, operation := range doc.Paths[path].Operations() {
			p := getRealPath(path)
			m := method
			for code, ref := range operation.Responses {
				// if p != "/store/order/:orderId" || code != "200" {
				// 	continue
				// }
				if p != "/" || code != "200" {
					continue
				}

				fmt.Println(p, code)
				responseSchema := core.NewResponseSchemaFromOpenApi3Ref(p, method, code, ref, "application/json")

				// engine must store for each path method code then check that
				d.expectationEngine.AddResponseSchemaForPathMethod(responseSchema)
				e := core.NewExample(responseSchema, d.GenerationConfig)
				if e == nil {
					fmt.Println("NOTHING")
					continue
				}

				status, _ := strconv.Atoi(code)
				d.expectationEngine.AddExpectation(expectation.PathMethod{
					Path:   p,
					Method: m,
				}, expectation.Expectation{
					Pririoty: 0,
					Response: expectation.Response{
						Body:   *e,
						Status: status,
					},
					Request: expectation.Request{
						Path:   p,
						Method: m,
					},
				})
			}
			d.EchoServer.Router().Add(m, p, func(ctx echo.Context) error {
				response, err := d.expectationEngine.GetResponseForRequest(p, m, ctx.Request())

				if err != nil {
					return ctx.JSON(500, GeneralError{
						Path:     ctx.Request().URL.Path,
						Method:   m,
						ErrorMsg: err.Error(),
					})
				}
				return ctx.JSON(response.Status, response.Body.Get())
			})
		}

	}
	d.expectationEngine.Load()

	end := time.Now()
	fmt.Println("TIME: ", end.Sub(start).Seconds())

	return d.EchoServer.Start(address)
}

func (d *Dolus) Start(address string) error {
	if !d.HideBanner {
		printBanner()
	}

	if d.expectationEngine == nil {
		generationConfig := d.GenerationConfig
		generationConfig.SetNonRequiredFields = true
		d.expectationEngine = engine.NewDolusExpectationEngine(generationConfig)
	}

	d.expectationEngine.AddExpectationsFromFiles(d.expectationFiles...)
	return d.startHttpServer(address)
}

func (d *Dolus) AddExpectations(files ...string) {
	d.expectationFiles = append(d.expectationFiles, files...)
}

func (d *Dolus) AddCueExpectationsFromFolder() {

}
