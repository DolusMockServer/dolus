package main

import (
	"encoding/json"
	"fmt"

	"github.com/MartinSimango/dolus"
	"github.com/MartinSimango/dstruct/generator"
)

type B struct {
	P string `json:"P"`
	L []int  `json:"L"`
	A *A
}
type A struct {
	B string `json:"B"`
	C *int   `json:"C"`
	F *B     `json:"F"`
	A *A
}

type L struct {
	AL string
	T  int
}

type X struct {
	Cx int
	Bx string
	Ax int
	Gx int
	Hx int
	Ix int
	Kx int
}

type Gx struct {
	Cx int
	Bx string
	Ax int
	Hx int
	Ix int
	Kx int
}

type AYA struct {
	X   *X `json:"x"`
	B   *any
	Int *int
}

func Print(strct any) string {
	val, _ := json.MarshalIndent(strct, "", "\t")
	return string(val)
}

func main() {

	d := dolus.New()
	d.AddExpectations("ideal.cue")
	d.GenerationConfig.
		SetValueGenerationType(generator.UseDefaults).
		SetNonRequiredFields(true)

	d.OpenAPIspec = "openapi-pet.yaml"
	if err := d.Start(fmt.Sprintf(":%d", 1080)); err != nil {
		fmt.Println(err)
	}

}
