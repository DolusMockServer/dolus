package task

import "github.com/MartinSimango/dstruct/generator"

func RegisterDolusTasks() {
	generator.AddTask(&GenInt32Task{})
}
