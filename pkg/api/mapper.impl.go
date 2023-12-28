package api

import (
	"fmt"

	"github.com/DolusMockServer/dolus/internal/server"
	"github.com/DolusMockServer/dolus/pkg/expectation/cue"
)

type MapperImpl struct{}

var _ Mapper = &MapperImpl{}

func NewMapper() *MapperImpl {
	return &MapperImpl{}
}

func (mi *MapperImpl) MapCueExpectations(
	expectations []cue.Expectation,
) ([]server.Expectation, error) {
	var apiServerExpectations []server.Expectation
	for _, cueExpectation := range expectations {
		apiServerExpectation, err := cueExpectationToApiExpectation(cueExpectation)
		if err != nil {
			return nil, err
		}
		apiServerExpectations = append(apiServerExpectations, *apiServerExpectation)
	}
	return apiServerExpectations, nil
}

func anyToMapOfKeyStringValueAny(a any) (*map[string]any, error) {
	if a != nil {
		if r, ok := a.(map[string]interface{}); ok {
			return &r, nil
		} else {
			return nil, fmt.Errorf("failed to convert %v to map[string]interface{}", a)
		}
	}
	return nil, nil
}

func callbackToApiCallback(callback *cue.Callback) (*server.Callback, error) {
	if callback != nil {
		callbackRequestBody, err := anyToMapOfKeyStringValueAny(callback.Request)
		if err != nil {
			return nil, err
		}
		return &server.Callback{
			HttpMethod:  string(callback.Method),
			RequestBody: callbackRequestBody,
			Timeout:     callback.Timeout,
			Url:         string(callback.Url),
		}, nil
	}
	return nil, nil
}

func cueExpectationToApiExpectation(cueExpectation cue.Expectation) (*server.Expectation, error) {
	requestBody, responseBody, err := getRequestAndResponseBody(cueExpectation)
	if err != nil {
		return nil, err
	}
	callback, err := callbackToApiCallback(cueExpectation.Callback)
	if err != nil {
		return nil, err
	}

	return &server.Expectation{
		Priority: cueExpectation.Priority,
		Request: server.Request{
			Method: string(cueExpectation.Request.Method),
			Path:   cueExpectation.Request.Path,
			Body:   requestBody,
		},
		Response: server.Response{
			Body:   responseBody,
			Status: cueExpectation.Response.Status,
		},
		Callback: callback,
	}, nil
}

// ExpectationToApiExpectation implements DolusApiFactory.
func getRequestAndResponseBody(
	expectation cue.Expectation,
) (*map[string]any, *map[string]any, error) {
	requestBody, err := anyToMapOfKeyStringValueAny(expectation.Request.Body)
	if err != nil {
		return nil, nil, err
	}
	responseBody, err := anyToMapOfKeyStringValueAny(expectation.Response.Body)
	if err != nil {
		return nil, nil, err
	}
	return requestBody, responseBody, nil
}
