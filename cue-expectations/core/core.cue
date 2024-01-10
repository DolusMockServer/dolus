// package core

// #Response: {
//   body: _
//   status: int
// }

// httpMethod: "GET" | "POST" | "HEAD" | "PUT" | "OPTIONS" | "TRACE" | "DELETE"

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


// #Request: {
//   path: string
//   method: httpMethod
//   body: _
//   params?: {
//     path?: [string]: #Matcher
//     query?: [string]: #Matcher
//   }
//   headers?: _
//   cookies?: _
// }

// #Matcher: { 
//   match: "eq" | "has" | "regex" | "not"
//   value: _
// }

// httpUrlRegex: =~"^(https?://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9-._~:/?#@$&'()*+,;=]*)?)$"

// #Callback:  {
//   request: _
//   // how long to wait
//   timeout: int | *1000
//   url: httpUrlRegex
//   method: httpMethod
// }


// #Expectation:  {
//     priority: int | *0
//     request: #Request
//     response:  #Response
//     callback?:  #Callback

// }



// #Expectations: {
//   expectations:  [...#Expectation]
// }



