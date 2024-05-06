package api

import (
	"fmt"

	"github.com/DolusMockServer/dolus/pkg/expectation"
)

// TOOD: see if there is a go package that can make mapping less verbose

type MapperImpl struct{}

var _ Mapper = &MapperImpl{}

func NewMapper() *MapperImpl {
	return &MapperImpl{}
}

func (mi *MapperImpl) MapToApiExpectations(
	expectations []expectation.Expectation,
) ([]Expectation, error) {
	var apiExpectations []Expectation
	for _, expct := range expectations {
		apiServerExpectation, err := expectationToApiExpectation(expct)
		if err != nil {
			return nil, err
		}
		apiExpectations = append(apiExpectations, *apiServerExpectation)
	}
	return apiExpectations, nil
}

func (mi *MapperImpl) MapToApiExpectation(expct expectation.Expectation) (*Expectation, error) {
	return expectationToApiExpectation(expct)
}

func (mi *MapperImpl) MapToExpectation(expct Expectation) (*expectation.Expectation, error) {
	callback, err := apiCallbackToCallback(expct.Callback)
	if err != nil {
		return nil, err
	}

	return &expectation.Expectation{
		Priority: expct.Priority,
		Request: expectation.Request{
			Method: expct.Request.Method,
			Path:   expct.Request.Path,
			Body:   expct.Request.Body,
		},
		Response: expectation.Response{
			Body:   expct.Response.Body,
			Status: expct.Response.Status,
		},
		Callback: callback,
	}, nil
}

func expectationToApiExpectation(expct expectation.Expectation) (*Expectation, error) {
	requestBody, responseBody, err := getRequestAndResponseBody(expct)
	if err != nil {
		return nil, err
	}
	callback, err := callbackToApiCallback(expct.Callback)
	if err != nil {
		return nil, err
	}

	return &Expectation{
		Priority: expct.Priority,
		Request: Request{
			Method: string(expct.Request.Method),
			Path:   expct.Request.Path,
			Body:   requestBody,
		},
		Response: Response{
			Body:   responseBody,
			Status: expct.Response.Status,
		},
		Callback: callback,
	}, nil
}

func getRequestAndResponseBody(
	expectation expectation.Expectation,
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

func anyToMapOfKeyStringValueAny(a any) (*map[string]any, error) {
	if a != nil {
		if r, ok := a.(map[string]interface{}); ok {
			return &r, nil
		} else {
			// TODO: see if this scenario is possible
			return nil, fmt.Errorf("failed to convert %v to map[string]interface{}", a)
		}
	}
	return nil, nil
}

func apiCallbackToCallback(apiCallback *Callback) (*expectation.Callback, error) {
	if apiCallback != nil {
		callbackRequestBody, err := anyToMapOfKeyStringValueAny(apiCallback.RequestBody)
		if err != nil {
			return nil, err
		}
		return &expectation.Callback{
			Method:  apiCallback.HttpMethod,
			Request: callbackRequestBody,
			Timeout: apiCallback.Timeout,
			Url:     apiCallback.Url,
		}, nil
	}
	return nil, nil
}

func callbackToApiCallback(callback *expectation.Callback) (*Callback, error) {
	if callback != nil {
		callbackRequestBody, err := anyToMapOfKeyStringValueAny(callback.Request)
		if err != nil {
			return nil, err
		}
		return &Callback{
			HttpMethod:  string(callback.Method),
			RequestBody: callbackRequestBody,
			Timeout:     callback.Timeout,
			Url:         string(callback.Url),
		}, nil
	}
	return nil, nil
}
