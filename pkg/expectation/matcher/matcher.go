package matcher

type Matcher[T any, M any] interface {
	Matches(value *M) bool
	GetValue() *T
}

type CueMatcher[T any] struct {
	MatchExpression string `json:"match"`
	Value           *T     `json:"value"` // pointers allows a match value to be missing
}

func (m CueMatcher[T]) GetValue() *T {
	return m.Value
}
