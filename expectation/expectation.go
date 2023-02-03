package expectation

import "github.com/MartinSimango/dolus/core"

type PathMethod struct {
	Path   string
	Method string
}

type PathMethodStatus struct {
	PathMethod
	Status string
}

type Response struct {
	Body   core.Example
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
	Pririoty int
	Response Response
	Request  Request
	// RequestMatcher
}

// import (
// 	"net/http"

// 	"github.com/MartinSimango/dolus/pkg/example"
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
