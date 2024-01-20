package expectation

import (
	"fmt"
)

type ExpectationFieldError struct {
	FieldName string
	Err       error
}

type ExpectationError struct {
	Expectation
	ExpectationFieldErrors []ExpectationFieldError
}

func (e ExpectationError) Error() string {
	retString := fmt.Sprintf(
		"Expectation Error:\nPath: %s\nMethod: %s\nPriority: %d\nErrors:\n",
		e.Expectation.Request.Path,
		e.Expectation.Request.Method,
		e.Expectation.Priority,
	)
	for i, err := range e.ExpectationFieldErrors {
		retString = fmt.Sprintf("%s%d. %+v\n", retString, (i + 1), err)
	}
	return retString
}
