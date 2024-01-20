package schema

import (
	"net/url"
)

type RequestParameters struct {
	PathParams  map[string]string
	QueryParams url.Values
}
