package cue

import (
	"time"
)

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

// struct taking from net/http package
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

type Response struct {
	Body    *any                 `json:"body"`
	Status  int                  `json:"status"`
	Headers *map[string][]string `json:"headers"`
	Cookies *[]Cookie            `json:"cookies"`
}

type Request struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Body   *any   `json:"body"`
	Params *struct {
		Path  *map[string]Matcher `json:"path"`
		Query *map[string]Matcher `json:"query"`
	} `json:"params"`

	Headers *map[string][]string `json:"headers"`
	Cookies *[]Cookie            `json:"cookies"`
}

type Matcher struct {
	Match string `json:"match"`
	Value *any   `json:"value"`
}

type Callback struct {
	Request any    `json:"request"`
	Timeout int    `json:"timeout"`
	Url     string `json:"url"`
	Method  string `json:"method"`
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
