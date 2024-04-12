package expectation

import (
	"fmt"
	"net/url"
	"strings"

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

func (e *Expectation) AddRequestParameterMatchers(pathParams map[string]string, queryParams url.Values) error {
	if e.Request.Parameters == nil {
		e.Request.Parameters = &RequestParameters{}
	}
	e.addQueryParameters(queryParams)

	return e.addPathParameters(pathParams)

}

func (e *Expectation) addPathParameters(pathParams map[string]string) error {
	if e.Request.Parameters.Path == nil {
		e.Request.Parameters.Path = make(map[string]any)
	}
	for k, v := range pathParams {
		matchType := "eq"
		value := v
		if v == ":"+k {
			matchType = "has"
		} else if strings.TrimSpace(v) == "" {
			return fmt.Errorf("path parameter '%s' is empty", k)
		}
		e.Request.Parameters.Path[k] = NewStringMatcher(value, matchType)

	}
	return nil
}

func (e *Expectation) addQueryParameters(queryParams url.Values) {
	if e.Request.Parameters.Query == nil {
		e.Request.Parameters.Query = make(map[string]any)
	}
	for k, v := range queryParams {
		value := v
		e.Request.Parameters.Query[k] = NewStringArrayMatcher(value, "eq")
	}
}

// ValidateRequestParameters validates the request parameters against the schema
func (e *Expectation) ValidateRequestParameters(requestParamProp schema.RequestParameterProperty) error {

	// Validate Path and Query Parameters
	if err := checkParametersExistence(PATH_PARAM, requestParamProp.PathParameterProperties, e.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkRequiredParameters[*StringMatcher](PATH_PARAM, requestParamProp.PathParameterProperties, e.Request.Parameters.Path); err != nil {
		return err
	}

	if err := checkParametersExistence(QUERY_PARAM, requestParamProp.QueryParameterProperties, e.Request.Parameters.Query); err != nil {
		return err
	}

	if err := checkRequiredParameters[*StringArrayMatcher](QUERY_PARAM, requestParamProp.QueryParameterProperties, e.Request.Parameters.Query); err != nil {
		return err
	}
	return nil
}

func checkParametersExistence(paramType string, properties schema.ParameterProperties, parameters map[string]any) error {
	for name := range parameters {
		if properties[name] == nil {
			return fmt.Errorf("%s parameter '%s' does not exist", paramType, name)
		}
	}
	return nil
}

func checkRequiredParameters[T Matcher](paramType string, properties schema.ParameterProperties, values map[string]any) error {
	for value, param := range properties {

		if param.Required && (values[value] == nil || values[value].(T).GetValue() == nil) {
			return fmt.Errorf("required %s parameter '%s' is missing", paramType, value)
		}
	}
	return nil
}
