package dolus

import (
	"fmt"

	"github.com/DolusMockServer/dolus/server"
	"github.com/MartinSimango/dolus-expectations/pkg/dolus"
)

type DolusApiFactory interface {
	RawCueExpectationsToApiExpectations(expectations []dolus.Expectation) ([]server.Expectation, error)
}

type DolusApiFactoryImpl struct{}

var _ DolusApiFactory = &DolusApiFactoryImpl{}

func NewDolusApiFactoryImpl() *DolusApiFactoryImpl {
	return &DolusApiFactoryImpl{}
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

func callbackToApiCallback(callback *dolus.Callback) (*server.Callback, error) {
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

func rawCueExpectationToApiExpectation(e dolus.Expectation) (*server.Expectation, error) {
	requestBody, responseBody, err := getRequestAndResponseBody(e)
	if err != nil {
		return nil, err
	}
	callback, err := callbackToApiCallback(e.Callback)
	if err != nil {
		return nil, err
	}

	return &server.Expectation{
		Priority: e.Priority,
		Request: server.Request{
			Method: string(e.Request.Method),
			Path:   e.Request.Path,
			Body:   requestBody,
		},
		Response: server.Response{
			Body:   responseBody,
			Status: e.Response.Status,
		},
		Callback: callback,
	}, nil

}

// ExpectationToApiExpectation implements DolusApiFactory.
func getRequestAndResponseBody(expectation dolus.Expectation) (*map[string]any, *map[string]any, error) {
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

func (daf *DolusApiFactoryImpl) RawCueExpectationsToApiExpectations(expectations []dolus.Expectation) ([]server.Expectation, error) {
	var apiServerExpectations []server.Expectation
	for _, e := range expectations {
		apiServerExpectation, err := rawCueExpectationToApiExpectation(e)
		if err != nil {
			return nil, err
		}
		apiServerExpectations = append(apiServerExpectations, *apiServerExpectation)
	}
	return apiServerExpectations, nil

}
