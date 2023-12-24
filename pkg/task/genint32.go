package task

import (
	"github.com/MartinSimango/dstruct/generator"
)

const (
	// TODO: just make this a string
	GenInt32 generator.TaskName = "GenInt32"
)

type GenInt32Task struct{}

var _ generator.Task = &GenInt32Task{}

// Tags implements Task.
func (g *GenInt32Task) Name() string {
	return string(GenInt32)
}

// GenerationFunction implements Task.
func (g *GenInt32Task) GenerationFunction(
	taskProperties generator.TaskProperties,
) generator.GenerationFunction {
	generator.ValidateParamCount(g, taskProperties)
	params := getNumberParams[int32](
		taskProperties.FieldName,
		taskProperties.Parameters,
		string(GenInt32),
	)
	return generator.GenerateNumberFunc(params.min, params.max)
}

func (g *GenInt32Task) ExpectedParameterCount() int {
	return 2
}
