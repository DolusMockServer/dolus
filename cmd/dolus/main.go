package main

import (
	"fmt"

	"github.com/MartinSimango/dolus/pkg/dolus"
	"github.com/MartinSimango/dolus/pkg/example"
	"github.com/MartinSimango/dolus/pkg/generator"
	"github.com/MartinSimango/dolus/pkg/schema"
)

func main() {

	d := dolus.New()
	d.ExpectationEngine.Expectations = []string{"ideal.cue"}
	d.ExpectationEngine.Start()
	s := schema.NewSchemaFromCueValue("/store/order/:orderId", "GET", "200", d.ExpectationEngine.Schemas[0])
	c := generator.NewGenerationConfig()
	e := example.New(s, *c)
	fc := e.GetFieldGenerationConfig("body.test")
	if fc != nil {
		fc.SetNonRequiredFields = true
		fc.Int64Min = 4
		fc.Int64Max = 4
		fc.Float64Min = 2
		fc.Float64Max = 4
		// fc.ValueGenerationType = generator.GenerateOnce
	}
	d.ResponseRepository.Add("/store/order/:orderId", "GET", "200", e)

	d.OpenAPIspec = "openapi-pet.yaml"
	if err := d.Start(fmt.Sprintf(":%d", 1080)); err != nil {
		fmt.Println(err)
	}

}
