package task

import (
	"github.com/MartinSimango/dstruct/generator"
)

func RegisterDolusTasks() {
	generator.AddTask(&GenInt32Task{})
	generator.AddTask(&GenIntTask{})
}

func RegisterTask(task generator.Task) error {
	return generator.AddTask(task)
}
