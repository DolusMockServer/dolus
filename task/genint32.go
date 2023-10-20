package task

import (
	"fmt"
	"strconv"

	"github.com/MartinSimango/dstruct/generator"
)

const (
	GenInt32 generator.TaskName = "GenInt32"
)

type genInt32Params struct {
	min int32
	max int32
}

type GenInt32Task struct{}

var _ generator.Task = &GenInt32Task{}

// Tags implements Task.
func (g *GenInt32Task) Name() string {
	return string(GenInt32)
}

// GenerationFunction implements Task.
func (g *GenInt32Task) GenerationFunction(taskProperties generator.TaskProperties) generator.GenerationFunction {
	generator.ValidateParamCount(g, taskProperties)
	params := g.getInt32Params(taskProperties.FieldName, taskProperties.Parameters)
	return generator.GenerateNumberFunc(params.min, params.max)

}

func (g *GenInt32Task) ExpectedParameterCount() int {
	return 2
}

func (g *GenInt32Task) getInt32Params(fieldName string, params []string) genInt32Params {

	param_1, err := strconv.Atoi(params[0])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))
	}

	param_2, err := strconv.Atoi(params[1])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))
	}

	if param_1 > param_2 {
		err = fmt.Errorf("min must be less or equal to the max value min = %d max = %d", param_1, param_2)
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))

	}

	return genInt32Params{
		min: int32(param_1),
		max: int32(param_2),
	}
}
