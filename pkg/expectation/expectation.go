package expectation

import (
	"fmt"

	"github.com/MartinSimango/dstruct"

	"github.com/DolusMockServer/dolus-expectations/pkg/dolus"
)

type PathMethod struct {
	OpenApiPath string
	Path        string
	Method      string
}

type PathMethodStatus struct {
	PathMethod
	Status string
}

func PathMethodFromExpectation(expectation DolusExpectation) PathMethod {
	return PathMethod{
		OpenApiPath: expectation.Request.OpenApiPath,
		Method:      expectation.Request.Method,
	}
}

func PathMethodStatusExpectation(expectation DolusExpectation) PathMethodStatus {
	return PathMethodStatus{
		PathMethod: PathMethodFromExpectation(expectation),
		Status:     fmt.Sprintf("%d", expectation.Response.Status),
	}
}

type DolusResponse struct {
	Body   dstruct.GeneratedStruct
	Status int
	// headers
}

type DolusRequest struct {
	OpenApiPath string
	Method      string
	Body        any
	// headers
}

type DolusExpectation struct {
	Priority          int
	Response          DolusResponse
	Request           DolusRequest
	RawCueExpectation *dolus.Expectation
	// RequestMatcher
}

type ExpectationFieldError struct {
	FieldName string
	Err       error
}

type ExpectationError struct {
	DolusExpectation
	ExpectationFieldErrors []ExpectationFieldError
}

func (e ExpectationError) Error() string {
	retString := fmt.Sprintf(
		"Expectation Error:\nPath: %s\nMethod: %s\nPriority: %d\nErrors:\n",
		e.DolusExpectation.Request.OpenApiPath,
		e.DolusExpectation.Request.Method,
		e.DolusExpectation.Priority,
	)
	for i, err := range e.ExpectationFieldErrors {
		retString = fmt.Sprintf("%s%d. %+v\n", retString, (i + 1), err)
	}
	return retString
}
