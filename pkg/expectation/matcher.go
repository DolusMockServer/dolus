package expectation

type CueMatcher[T any] struct {
	MatchExpression string `json:"match"`
	Value           *T     `json:"value"` // pointers allows a match value to be missing
}

func (m CueMatcher[T]) GetValue() *T {
	return m.Value
}
