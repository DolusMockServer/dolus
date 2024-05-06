package expectation

type ExpectationType int

const (
	Default ExpectationType = iota
	Custom
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
	// MatchRules MatchRules `json:"-"`
}

// type Rule[T any] struct {
// 	MatchType string `json:"matchType"`
// 	Value     T      `json:"value"`
// }

// type MatchRules struct {
// 	Headers         map[string]Rule[[]string]
// 	PathParameters  map[string]Rule[string]
// 	QueryParameters map[string]Rule[[]string]
// 	Cookies         map[string]Rule[Cookie]
// }

/*
	MatchRules: {
		"headers": {
			"Content-Type": {
				"matchType": "eq",
				 "value": "application/json"
			}
		}
		"PathParameters": {
			"orderId": {
				"matchType": "eq",
				 "value": "2"
		}

	}
*/
