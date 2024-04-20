package matcher

type StringMatcher struct {
	SimpleMatcher[string]
}

var _ Matcher[string] = &StringMatcher{}

func NewStringMatcher(value, matchType string) *StringMatcher {
	return &StringMatcher{
		SimpleMatcher: SimpleMatcher[string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}

}

func (m StringMatcher) Matches(value string) bool {
	switch m.MatchExpression {
	case "eq":
		return m.Value == value
	case "has":
		return true
	}
	return false
}
