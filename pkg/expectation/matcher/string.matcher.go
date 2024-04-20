package matcher

type StringMatcher struct {
	CueMatcher[string]
}

var _ Matcher[string, string] = &StringMatcher{}

func NewStringMatcher(value *string, matchType string) *StringMatcher {
	return &StringMatcher{
		CueMatcher: CueMatcher[string]{
			MatchExpression: matchType,
			Value:           value,
		},
	}

}

func (m StringMatcher) Matches(value *string) bool {
	switch m.MatchExpression {
	case "eq":
		return *m.Value == *value
	case "has":
		return value != nil
	}
	return false
}
