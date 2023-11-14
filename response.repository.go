package dolus

type GeneralError struct {
	Path     string
	Method   string
	ErrorMsg string
}

// import (
// 	"github.com/DolusMockServer/dolus/pkg/example"
// 	"github.com/DolusMockServer/dolus/pkg/expectation"
// 	"github.com/labstack/echo/v4"
// )

// type ExpectationRepository struct {
// 	Expectations map[expectation.PathMethod][]expectation.Expectation
// }

// func NewResponseRepository() *ExpectationRepository {
// 	return &ExpectationRepository{
// 		Expectations: make(map[expectation.PathMethod][]expectation.Expectation),
// 	}
// }

// func (repo *ExpectationRepository) GetEchoResponse(path, method string, ctx echo.Context) error {
// 	// need to look at request path, method
// 	// look at path and method - see if any expectations are there for

// 	// look for one with highest priority - wth path and method of which matches the request
// 	// -- of those look which one matches request if anything to match
// 	// -- look at the expectation examples and return

// 	for _, expectation := range repo.Expectations[PathMethod{Path: path, Method: method}] {

// 		if expectation.StatusCode == "200" {
// 			return ctx.JSON(200, expectation.Example.Get())

// 		}
// 	}
// 	return ctx.JSON(500, GeneralError{
// 		Path:     path,
// 		Method:   method,
// 		ErrorMsg: "No expectation found for path and HTTP method.",
// 	})

// }

// func (repo *ExpectationRepository) Add(path, method, code string, example *example.Example) {
// 	if example == nil {
// 		return
// 	}
// 	pathMethod := PathMethod{
// 		Path:   path,
// 		Method: method,
// 	}

// 	repo.Expectations[pathMethod] = append(repo.Expectations[pathMethod],
// 		B{
// 			Priority:   0,
// 			StatusCode: code,
// 			Request:    nil,
// 			Example:    example})

// }

// // func (repo *ResponseRepository) GetResponse(operation, path string, ctx echo.Context) error {

// // }
