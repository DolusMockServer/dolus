package expectation

import (
	"fmt"

	"github.com/MartinSimango/dolus-expectations/pkg/dolus"
	"github.com/MartinSimango/dstruct"
)

type PathMethod struct {
	Path   string
	Method string
}

type PathMethodStatus struct {
	PathMethod
	Status string
}

func PathMethodFromExpectation(expectation Expectation) PathMethod {
	return PathMethod{
		Path:   expectation.Request.Path,
		Method: expectation.Request.Method,
	}
}

func PathMethodStatusExpectation(expectation Expectation) PathMethodStatus {
	return PathMethodStatus{
		PathMethod: PathMethodFromExpectation(expectation),
		Status:     fmt.Sprintf("%d", expectation.Response.Status),
	}
}

type Response struct {
	Body   dstruct.GeneratedStruct
	Status int
	//headers
}

type Request struct {
	Path   string
	Method string
	Body   any
	//headers

}

type Expectation struct {
	Priority          int
	Response          Response
	Request           Request
	RawCueExpectation *dolus.Expectation
	// RequestMatcher
}

type ExpectationFieldError struct {
	FieldName string
	Err       error
}

type ExpectationError struct {
	Expectation
	ExpectationFieldErrors []ExpectationFieldError
}

func (e ExpectationError) Error() string {
	retString := fmt.Sprintf("Expectation Error:\nPath: %s\nMethod: %s\nPriority: %d\nErrors:\n", e.Expectation.Request.Path,
		e.Expectation.Request.Method,
		e.Expectation.Priority)
	for i, err := range e.ExpectationFieldErrors {
		retString = fmt.Sprintf("%s%d. %+v\n", retString, (i + 1), err)
	}
	return retString
}

// import (
// 	"net/http"

// 	"github.com/DolusMockServer/dolus/pkg/example"
// )

// // type NoExpectationType
// type ExpectationResponse struct {
// 	Body   example.Example
// 	Status int
// }

// type Expectation struct {
// 	Priority int
// 	Path     string
// 	Request  http.Request
// 	Response http.Response
// 	Example  example.Example
// }

// // Path     string   `json:"path"`
// // Method   string   `json:"method"`
// // Pririoty int      `json:"priority"`
// // Response Response `json:"response"`
// // Request  Request  `json:"request"`

// // TODO have different types of expe

// // TODO server should have cap of when to send 429 if certain number of requests are coming through

// // If No expectation
// //   -- check if operation has example
// //      return example
// //   -- check if operation has schema
// //      	check no expectation type (GENERATED or USED Default type values)
// //   --
// // -- if no schema return error (internal server with message saying response could not be given)
// // -- Check type of expectation (GENERATERD tpyes)

// // -- return 200 for any request for any Operation if no schema

// // Request Operation Path -> 200
