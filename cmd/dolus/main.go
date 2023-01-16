package main

import (
	"github.com/MartinSimango/dolus/pkg/dolus"
)

func main() {

	// // str := string(b) // convert content to a 'string'
	// load.Instances()

	// ctx := cuecontext.New()
	// i := cuebuilder.NewContext().NewInstance("expectations", nil)
	// i.Dependencies()
	// i.AddFile("expectations/ideal.cue", nil)
	// id := cuebuilder.NewContext().NewInstance("dolus", nil)
	// id.AddFile("expectations/dolus/expectation.cue", nil)
	// id.AddFile("expectations/dolus/int.cue", nil)

	// // i.AddFile("expectations/dolus/expectation.cue", nil)

	// // file, err := cueparser.ParseFile("expectations/ideal.cue", nil, cueparser.ParseComments)
	// // if err != nil {
	// // 	fmt.Print(err)
	// // 	os.Exit(1)
	// // }

	// instance, _ := ctx.BuildInstances([]*cuebuilder.Instance{i, id})
	// fmt.Println("GHL: ", instance[1])
	// // b, _ := instance.MarshalJSON()
	// fmt.Println("JS: ", i.Dependencies())
	// ctx.BuildFile(file)

	// it, _ := instance.List()

	// fmt.Println("begin: ", it.Value())

	// fmt.Println("a: ", it.Next())
	// fmt.Println("v: ", it.Value())
	// b, _ := it.Value().MarshalJSON()
	// var g G
	// err = json.Unmarshal(b, &g)
	// fmt.Println(err)
	// fmt.Println(g.Response)

	// m := (g.Response).(map[string]interface{})
	// be, _ := schema.BuildExample(m, "", "")
	// nv := reflect.New(reflect.ValueOf(be).Type()).Interface()
	// dstruct.New(nv).Print()

	// ds := dstruct.New(&v)
	// ds.Print()
	// fmt.Println("A: ", ds.GetField("path"))

	// fmt.Println(it.Value().Struct())

	// fmt.Println("a: ", it.Next())
	// fmt.Println("v: ", it.Value())

	d := dolus.New()
	d.ExpectationEngine.Expectations = []string{"ideal.cue"}
	d.ExpectationEngine.Start()
	// d.OpenAPIspec = "openapi-pet.yaml"
	// if err := d.Start(fmt.Sprintf(":%d", 1080)); err != nil {
	// 	fmt.Println(err)
	// }

}
