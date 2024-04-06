package expectation


// Add contraints
import "time"

httpMethod: "GET" | "POST" | "HEAD" | "PUT" | "OPTIONS" | "TRACE" | "DELETE"

httpUrlRegex: =~"^(https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9-._~:/?#@$&'()*+,;=]*)?)$"

#cookieValue : #CookieMatcher | #Cookie

#Request: #Request & {
    method: httpMethod
    headers: {[string]: {#HeaderMatcher | #HeaderValueType } }
    cookies: [...#cookieValue]      // {[string]: {#CookieMatcher | #Cookie } }
}


#PathValueType: string | int | float | bool 
#QueryValueType: [...string] | [...int] | int | string | float | bool
#HeaderValueType: [...string] | string | int 

#RequestParameters: #RequestParameters & {
    path: {[string]: {#PathMatcher | #PathValueType} }
    query: {[string]: {#QueryMatcher | #QueryValueType} }
}


#HeaderMatcher: #Matcher & {
    match: string| *"eq"
} & (
    { match: "has", value: null} |
    { match: "eq" | "regex" | "not", value: #HeaderValueType} )

#PathMatcher: #Matcher & {
    match: string| *"eq"
} & (
    { match: "has", value: null} |
    { match: "eq" | "regex" | "not", value: #PathValueType } )


#QueryMatcher: #Matcher & {
    match: string| *"eq"
} & (
    { match: "has", value: null} |
    { match: "eq" | "regex" | "not", value: #QueryValueType } )

#CookieMatcher: #Matcher & {
    match: string| *"eq"
} & (
    { match: "has", value: null} |
    { match: "eq" , value: #Cookie })



#Callback: #Callback & {
    timeout: int | *1000
    url: httpUrlRegex
    method: httpMethod
}

#Expectation: #Expectation & {
    priority: int | *0
}


#Cookie: #Cookie & {
	value:      string | *"",             // Default value is an empty string
	path:       string | *"/",            // Default path is the root path
	domain:     string | *"",             // Default domain is empty (current domain)
	expires:    time.Time | *"0001-01-01T00:00:00Z",    // Default expiration time is the zero time
	rawExpires: string| *"",             // Default raw expiration is an empty string
	maxAge:   int | *0,                // Default MaxAge is 0 (no 'Max-Age' attribute specified)
	secure:   bool | *false,            // Default is not secure
	httpOnly: bool | *false,            // Default is not HttpOnly
	sameSite: #SameSite | *0,  // Default SameSite is SameSiteDefault
	raw:      string| *"",               // Default raw value is an empty string
	// Unparsed: ,              // Default unparsed attributes are nil
}

