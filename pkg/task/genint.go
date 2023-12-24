package task

import (
	"github.com/MartinSimango/dstruct/generator"
)

const (
	// TODO: just make this a string
	GenInt generator.TaskName = "GenInt"
)

type GenIntTask struct{}

var _ generator.Task = &GenIntTask{}

// Tags implements Task.
func (g *GenIntTask) Name() string {
	return string(GenInt)
}

// GenerationFunction implements Task.
func (g *GenIntTask) GenerationFunction(
	taskProperties generator.TaskProperties,
) generator.GenerationFunction {
	generator.ValidateParamCount(g, taskProperties)
	params := getNumberParams[int](
		taskProperties.FieldName,
		taskProperties.Parameters,
		string(GenInt),
	)
	return generator.GenerateNumberFunc(params.min, params.max)
}

func (g *GenIntTask) ExpectedParameterCount() int {
	return 2
}
