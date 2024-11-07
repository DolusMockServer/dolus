package task

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

func RegisterDolusTasks() {
	// Integer tasks
	numberRangeConfig := config.NewNumberRangeConfig()
	RegisterTask(NewGenIntTask(numberRangeConfig))
	RegisterTask(NewGenInt8Task(numberRangeConfig))
	RegisterTask(NewGenInt16Task(numberRangeConfig))
	RegisterTask(NewGenInt32Task(numberRangeConfig))
	RegisterTask(NewGenInt64Task(numberRangeConfig))
	RegisterTask(NewGenUIntTask(numberRangeConfig))
	RegisterTask(NewGenUInt8Task(numberRangeConfig))
	RegisterTask(NewGenUInt16Task(numberRangeConfig))
	RegisterTask(NewGenUInt32Task(numberRangeConfig))
	RegisterTask(NewGenUInt64Task(numberRangeConfig))
}

func RegisterTask(task generator.Task) error {
	return generator.AddTask(task)
}
