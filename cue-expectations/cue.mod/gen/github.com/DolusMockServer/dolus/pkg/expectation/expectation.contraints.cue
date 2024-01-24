package expectation


// Add contraints

httpMethod: "GET" | "POST" | "HEAD" | "PUT" | "OPTIONS" | "TRACE" | "DELETE"

httpUrlRegex: =~"^(https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9-._~:/?#@$&'()*+,;=]*)?)$"

#Request: #Request & {
    method: httpMethod
    headers: {[string]: {#HeaderMatcher | #HeaderValueType } }
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


#Callback: #Callback & {
    timeout: int | *1000
    url: httpUrlRegex
    method: httpMethod
}

#Expectation: #Expectation & {
    priority: int | *0
}
