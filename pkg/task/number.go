package task

import (
	"fmt"
	"reflect"
	"strconv"
)

// TODO: change number type to be a number and not just any
type genNumberParam[N any] struct {
	min N
	max N
}

// TODO: account for floats by checking type of N
func getNumberParams[N any](fieldName string, params []string, taskName string) genNumberParam[N] {
	param_1, err := strconv.Atoi(params[0])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, taskName, err))
	}

	param_2, err := strconv.Atoi(params[1])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, taskName, err))
	}

	if param_1 > param_2 {
		err = fmt.Errorf(
			"min must be less or equal to the max value min = %d max = %d",
			param_1,
			param_2,
		)
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, taskName, err))

	}

	v := any(*new(N))
	var min, max any
	switch v.(type) {
	case int:
		min, max = param_1, param_2
	case int8:
		min, max = int8(param_1), int8(param_2)
	case int16:
		min, max = int16(param_1), int16(param_2)
	case int32:
		min, max = int32(param_1), int32(param_2)
	case int64:
		min, max = int64(param_1), int64(param_2)
	case uint:
		min, max = uint(param_1), uint(param_2)
	case uint8:
		min, max = uint8(param_1), uint8(param_2)
	case uint16:
		min, max = uint16(param_1), uint16(param_2)
	case uint32:
		min, max = uint32(param_1), uint32(param_2)
	case uint64:
		min, max = uint64(param_1), uint64(param_2)
	case float32:
		min, max = float32(param_1), float32(param_2)
	case float64:
		min, max = float64(param_1), float64(param_2)
	case uintptr:
		min, max = uintptr(param_1), uintptr(param_2)
	default:
		panic(fmt.Sprintf("Type not supported for getNumberRange: %s", reflect.TypeOf(v)))
	}
	return genNumberParam[N]{
		min: min.(N),
		max: max.(N),
	}

}
