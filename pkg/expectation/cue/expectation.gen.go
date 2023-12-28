// DO NOT EDIT! Code was generated by cue2gostruct v0.0.3-alpha.

package cue

type Response struct {
	Body   any `json:"body"`
	Status int `json:"status"`
}

type httpMethod string

type Request struct {
	Path    string     `json:"path"`
	Method  httpMethod `json:"method"`
	Body    any        `json:"body"`
	Headers any        `json:"headers"`
	Cookies any        `json:"cookies"`
}

type httpUrlRegex string

type Callback struct {
	Request any          `json:"request"`
	Timeout int          `json:"timeout"`
	Url     httpUrlRegex `json:"url"`
	Method  httpMethod   `json:"method"`
}

type Expectation struct {
	Priority int       `json:"priority"`
	Request  Request   `json:"request"`
	Response Response  `json:"response"`
	Callback *Callback `json:"callback"`
}

type Expectations struct {
	Expectations []Expectation `json:"expectations"`
}
