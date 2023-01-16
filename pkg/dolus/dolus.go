package dolus

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MartinSimango/dolus/pkg/example"
	"github.com/MartinSimango/dolus/pkg/expectation"
	"github.com/MartinSimango/dolus/pkg/generator"
	"github.com/MartinSimango/dolus/pkg/schema"
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
	OpenAPIspec        string
	HideBanner         bool
	HidePort           bool
	EchoServer         *echo.Echo
	ResponseRepository *ResponseRepository
	ExpectationEngine  *expectation.ExpectationEngine
}

type OperationResponse struct {
	Operation string
	Response  http.Response
}

type Paths map[string][]OperationResponse

func New() *Dolus {
	return &Dolus{
		HideBanner:         false,
		HidePort:           false,
		OpenAPIspec:        "openapi.yaml",
		ResponseRepository: NewResponseRepository(),
		ExpectationEngine:  &expectation.ExpectationEngine{},
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
				fmt.Println(p, code)

				s := schema.New(path, method, code, ref, "application/json")
				generatorConfig := generator.NewGenerationConfig()
				generatorConfig.SetNonRequiredFields = true
				generatorConfig.ValueGenerationType = generator.Generate
				e := example.New(s, *generatorConfig)
				if e == nil {
					continue
				}
				fc := e.GetFieldGenerationConfig("quantity")
				if fc != nil {
					fc.SetNonRequiredFields = true
					fc.Int32Min = 4
					fc.Int32Max = 4
					fc.ValueGenerationType = generator.GenerateOnce
				}

				d.ResponseRepository.Add(p, m, code, e)
			}
			d.EchoServer.Router().Add(m, p, func(ctx echo.Context) error {
				return d.ResponseRepository.GetEchoResponse(p, m, ctx)
			})
		}

	}
	end := time.Now()
	fmt.Println("TIME: ", end.Sub(start).Seconds())
	return d.EchoServer.Start(address)
}

func (d *Dolus) Start(address string) error {
	if !d.HideBanner {
		printBanner()
	}
	return d.startHttpServer(address)
}

func (d *Dolus) AddExpectation() {
}
