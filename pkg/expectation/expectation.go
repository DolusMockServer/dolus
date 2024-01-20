package expectation

import (
	"fmt"
	"net/url"

	"github.com/DolusMockServer/dolus/pkg/schema"
)

type ExpectationType int

const (
	OpenAPI ExpectationType = iota
	Cue
)

// // request matcher - used to match which requests this expectation should be applied to
// // action - what action to take, actions include response, forward, callback and error
// // times (optional) - how many times the action should be taken
// // timeToLive (optional) - how long the expectation should stay active
// // priority (optional) - matching is ordered by priority (highest first) then creation (earliest first)
// // id (optional) - used for updating an existing expectation (i.e. when the id matches)

// // requestMatcher (create tickets)
// // method - property matcher - done
// // path - property matcher - done
// // path parameters - key to multiple values matcher - in progress
// // query string parameters - key to multiple values matcher - to do (partial query matches)
// // headers - key to multiple values matcher - to do
// // cookies - key to single value matcher - to do
// // body - body matchers - to do
// // secure - boolean value, true for HTTPS - to do (default to false)

type Expectations struct {
	Expectations []Expectation `json:"expectations"`
}
type Expectation struct {
	Priority int       `json:"priority"`
	Request  Request   `json:"request"`
	Response Response  `json:"response"`
	Callback *Callback `json:"callback"`
}

func (e *Expectation) AddRequestParameterMatchers(pathParams map[string]string, queryParams url.Values) {
	if e.Request.Parameters == nil {
		e.Request.Parameters = &RequestParameters{}
	}
	e.addPathParameters(pathParams)
	e.addQueryParameters(queryParams)

}

func (e *Expectation) addPathParameters(pathParams map[string]string) {
	if e.Request.Parameters.Path == nil {
		e.Request.Parameters.Path = make(map[string]Matcher[string])
	}
	for k, v := range pathParams {
		matchType := "eq"
		value := v
		if v == ":"+k {
			matchType = "has"
		}
		e.Request.Parameters.Path[k] = Matcher[string]{
			Match: matchType,
			Value: &value,
		}
	}
}

func (e *Expectation) addQueryParameters(queryParams url.Values) {
	if e.Request.Parameters.Query == nil {
		e.Request.Parameters.Query = make(map[string]Matcher[[]string])
	}
	for k, v := range queryParams {
		value := v
		e.Request.Parameters.Query[k] = Matcher[[]string]{
			Match: "eq",
			Value: &value,
		}
	}
}

// ValidateRequestParameters validates the request parameters against the schema
func (e *Expectation) ValidateRequestParameters(requestParamProp schema.RequestParameterProperty) error {

	// Validate Path and Query Parameters
	if err := checkParametersExistence("path", requestParamProp.PathParameterProperties, e.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkRequiredParameters("path", requestParamProp.PathParameterProperties, e.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkParametersExistence("query", requestParamProp.QueryParameterProperties, e.Request.Parameters.Query); err != nil {
		return err
	}

	if err := checkRequiredParameters("query", requestParamProp.QueryParameterProperties, e.Request.Parameters.Query); err != nil {
		return err
	}
	return nil
}

func checkParametersExistence[T any](paramType string, properties schema.ParameterProperties, parameters map[string]Matcher[T]) error {
	for name := range parameters {
		if properties[name] == nil {
			return fmt.Errorf("%s parameter '%s' does not exist", paramType, name)
		}
	}
	return nil
}

func checkRequiredParameters[T any](paramType string, properties schema.ParameterProperties, values map[string]Matcher[T]) error {
	for value, param := range properties {

		if param.Required && values[value].Value == nil {
			return fmt.Errorf("required %s parameter '%s' is missing", paramType, value)
		}
	}
	return nil
}
