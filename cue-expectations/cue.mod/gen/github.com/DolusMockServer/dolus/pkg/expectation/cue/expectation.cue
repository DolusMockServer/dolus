package cue


// Add contraints

httpMethod: "GET" | "POST" | "HEAD" | "PUT" | "OPTIONS" | "TRACE" | "DELETE"

httpUrlRegex: =~"^(https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9-._~:/?#@$&'()*+,;=]*)?)$"

#Request: #Request & {
    method: httpMethod
}

#Matcher: #Matcher & {
    match: string| *"eq"
} & (
    { match: "has", value: null} |
    { match: "eq" | "regex" | "not", value: string } )

#Callback: #Callback & {
    timeout: int | *1000
    url: httpUrlRegex
    method: httpMethod
}

#Expectation: #Expectation & {
    priority: int | *0
}