package task

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
	"github.com/MartinSimango/dstruct/generator/core"
)

const (
	GenInt    string = "GenInt"
	GenInt8   string = "GenInt8"
	GenInt16  string = "GenInt16"
	GenInt32  string = "GenInt32"
	GenInt64  string = "GenInt64"
	GenUInt   string = "GenUInt"
	GenUInt8  string = "GenUInt8"
	GenUInt16 string = "GenUInt16"
	GenUInt32 string = "GenUInt32"
	GenUInt64 string = "GenUInt64"
)

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type GenIntTask[T Integer] struct {
	generator.BaseTask
	config.NumberRangeConfig
}

type GenIntTaskInstance[T Integer] struct {
	GenIntTask[T]
	genFunc generator.GenerationFunction
}

var _ generator.Task = &GenIntTask[int]{}

var _ generator.TaskInstance = &GenIntTaskInstance[int]{}

func (g GenIntTask[T]) Instance(params ...string) generator.TaskInstance {
	gti := &GenIntTaskInstance[T]{
		GenIntTask: g,
	}
	gti.SetParameters(params...)
	gti.NumberRangeConfig = g.NumberRangeConfig.Copy()
	gti.genFunc = core.GenerateNumberFunc[T](gti.NumberRangeConfig)
	return gti
}

// GenerationFunction implements Task.
func (g *GenIntTaskInstance[T]) GenerationFunction() generator.GenerationFunction {
	return g.genFunc
}

func (g *GenIntTaskInstance[T]) SetParameters(params ...string) {
	g.ValidateParamCount(params...)
	setParameters[T](params, g.Name(), g.NumberRangeConfig)
}

func setParameters[N Integer](params []string, taskName string, cfg config.NumberRangeConfig) {
	param_1, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error with task '%s': %s", taskName, err))
	}

	param_2, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error with task '%s': %s", taskName, err))
	}

	if param_1 > param_2 {
		err = fmt.Errorf(
			"min must be less or equal to the max value min = %d max = %d",
			param_1,
			param_2,
		)
		panic(fmt.Sprintf("error with task '%s': %s", taskName, err))

	}

	switch v := (any(*new(N))).(type) {
	case int:
		cfg.Int().SetRange(int(param_1), int(param_2))
	case int8:
		cfg.Int8().SetRange(int8(param_1), int8(param_2))
	case int16:
		cfg.Int16().SetRange(int16(param_1), int16(param_2))
	case int32:
		cfg.Int32().SetRange(int32(param_1), int32(param_2))
	case int64:
		cfg.Int64().SetRange(int64(param_1), int64(param_2))
	case uint:
		cfg.UInt().SetRange(uint(param_1), uint(param_2))
	case uint8:
		cfg.UInt8().SetRange(uint8(param_1), uint8(param_2))
	case uint16:
		cfg.UInt16().SetRange(uint16(param_1), uint16(param_2))
	case uint32:
		cfg.UInt32().SetRange(uint32(param_1), uint32(param_2))
	case uint64:
		cfg.UInt64().SetRange(uint64(param_1), uint64(param_2))
	default:
		panic(fmt.Sprintf("Type '%s' not supported for task '%s'", reflect.TypeOf(v), taskName))
	}
}

func NewGenIntTask(cfg config.NumberRangeConfig) *GenIntTask[int] {
	return &GenIntTask[int]{
		BaseTask:          *generator.NewBaseTask(GenInt, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenInt8Task(cfg config.NumberRangeConfig) *GenIntTask[int8] {
	return &GenIntTask[int8]{
		BaseTask:          *generator.NewBaseTask(GenInt8, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenInt16Task(cfg config.NumberRangeConfig) *GenIntTask[int16] {
	return &GenIntTask[int16]{
		BaseTask:          *generator.NewBaseTask(GenInt16, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenInt32Task(cfg config.NumberRangeConfig) *GenIntTask[int32] {
	return &GenIntTask[int32]{
		BaseTask:          *generator.NewBaseTask(GenInt32, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenInt64Task(cfg config.NumberRangeConfig) *GenIntTask[int64] {
	return &GenIntTask[int64]{
		BaseTask:          *generator.NewBaseTask(GenInt64, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenUIntTask(cfg config.NumberRangeConfig) *GenIntTask[uint] {
	return &GenIntTask[uint]{
		BaseTask:          *generator.NewBaseTask(GenUInt, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenUInt8Task(cfg config.NumberRangeConfig) *GenIntTask[uint8] {
	return &GenIntTask[uint8]{
		BaseTask:          *generator.NewBaseTask(GenUInt8, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenUInt16Task(cfg config.NumberRangeConfig) *GenIntTask[uint16] {
	return &GenIntTask[uint16]{
		BaseTask:          *generator.NewBaseTask(GenUInt16, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenUInt32Task(cfg config.NumberRangeConfig) *GenIntTask[uint32] {
	return &GenIntTask[uint32]{
		BaseTask:          *generator.NewBaseTask(GenUInt32, 2),
		NumberRangeConfig: cfg,
	}
}

func NewGenUInt64Task(cfg config.NumberRangeConfig) *GenIntTask[uint64] {
	return &GenIntTask[uint64]{
		BaseTask:          *generator.NewBaseTask(GenUInt64, 2),
		NumberRangeConfig: cfg,
	}
}
