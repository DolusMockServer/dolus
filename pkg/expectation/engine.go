package expectation

import (
	"encoding/json"
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/MartinSimango/dolus/pkg/dstruct"
	"github.com/MartinSimango/dolus/pkg/schema"
)

type ExpectationEngine struct {
	Expectations []string
	cueValues    []cue.Value
	Schemas      []any
}

func (e *ExpectationEngine) Start() {
	ctx := cuecontext.New()
	entrypoints := e.Expectations
	bis := load.Instances(entrypoints, &load.Config{
		Dir: "expectations",
	})

	for _, bi := range bis {
		// check for errors on the  instance
		// these are typically parsing errors
		if bi.Err != nil {
			fmt.Println("Error during load:", bi.Err)
			continue
		}
		value := ctx.BuildInstance(bi)

		fmt.Println("DIR: ", bi.Dir)
		if value.Err() != nil {
			fmt.Println("Error during build:", value.Err())
			continue
		}

		// Validate the value
		err := value.Validate()
		if err != nil {
			fmt.Println("Error during validation:", err)
			continue
		}
		e.cueValues = append(e.cueValues, value)
		// Print the value
		e.addResponseSchema(value)
	}
}

func (e *ExpectationEngine) addResponseSchema(instance cue.Value) {

	expectations, _ := instance.Value().Lookup("expectations").List()
	var expectation Expectation

	for expectations.Next() {

		marshalledJson, err := expectations.Value().MarshalJSON()
		if err != nil {
			fmt.Println("Error with expectation: ", err)
			continue
		}
		err = json.Unmarshal(marshalledJson, &expectation)
		if err != nil {
			fmt.Println("Error with expectation: ", err)
			continue
		}
		// fmt.Println("E: ", expectation.Response)
		// cue.Selector
		// TODO change to lookup
		val := expectations.Value().Lookup("response")
		response := (expectation.Response).(map[string]interface{})
		responseSchema, _ := schema.BuildExample(response, "", "", &val)
		dstruct.New(reflect.New(reflect.ValueOf(responseSchema).Type()).Interface()).Print()
		e.Schemas = append(e.Schemas, reflect.New(reflect.ValueOf(responseSchema).Type()).Interface())
		// fmt.Println("A: ", ds.GetField("path"))

	}
}
