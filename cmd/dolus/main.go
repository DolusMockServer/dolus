package main

import (
	"encoding/json"
	"fmt"

	"github.com/MartinSimango/dolus/dstruct"
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
	b := dstruct.ExtendStruct(AYA{X: &X{Bx: "nice", Cx: 10}})

	c := b.RemoveField("X.Gx")

	bMod := b.Build()
	// t := &Gx{Bx: "H", Cx: 2}
	v, _ := bMod.Get("x")
	fmt.Printf("V: %+v\n", v)

	cMod := c.Build()
	cMod.Set("x.Kx", 4)
	fmt.Println(cMod.Get("x.Kx"))

	fmt.Println(bMod.Get("x.Kx"))

	fmt.Println(Print(bMod.Instance()))

	// bBuild := b.Build()
	// bcBuild := bBuild

	// b.

	// // fmt.Println("T: ", reflect.TypeOf(bi))
	// // i := 3
	// bi := bcBuild.New()

	// t.Bx = "COOLIO"

	// fmt.Println(Print(bi))

	// bBuild.Set("x.B", "Hello")

	// v, _ := bBuild.Get("x.Bx")
	// fmt.Printf("B: %+v\n", v)
	// bi := bBuild.Instance()

	// c := dstruct.ExtendStruct(bi).Build().Instance()

	// fmt.Printf("C: %+v\n", c)

	// fmt.Printf("B: %+v\n", bi)

	// d := dolus.New()
	// d.AddExpectations("ideal.cue")
	// d.GenerationConfig.SetNonRequiredFields = true
	// // d.OpenAPIspec = "openapi-pet.yaml"
	// if err := d.Start(fmt.Sprintf(":%d", 1080)); err != nil {
	// 	fmt.Println(err)
	// }
	// // b := dstruct.NewStruct()

	// // b.AddField("P", "hi", `json:"P"`)
	// // b.AddField("L", []int{}, `json:"L"`)

	// d.AddField("Hello", "hi", `json:"Hello"`)
	// b := &B{
	// 	P: "pass",
	// 	L: []int{1, 2, 3},
	// }
	// a := A{
	// 	B: "cool",
	// 	// C: new(int),
	// 	F: b,
	// 	// A: &A{},
	// }

	// // fmt.Println(reflect.TypeOf(dstruct.ExtendStruct(a).RemoveField("F.P").Build().New()))
	// // a = nil

	// d.AddField("Guest", &a, `json:"Guest"`)

	// // // d.AddField("Guest", nil, `json:"Guest"`)

	// // // d.RemoveField("Hello")
	// // // d.RemoveField("Guest")

	// // // d.GetField()

	// ds := d.Build()
	// f, err := ds.GetField("Guest.B")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("Before: %+v\n", f)

	// ds.SetField("Guest.B", "show")
	// f, _ = ds.GetField("Guest.B")
	// fmt.Printf("After: %+v\n", f)

	// // fmt.Println(ds.New())
	// // m, err := json.MarshalIndent(ds.New(), "", "\t")
	// // if err != nil {
	// // 	fmt.Println("ERR")
	// // } else {
	// // 	fmt.Println(string(m))
	// // }

}
