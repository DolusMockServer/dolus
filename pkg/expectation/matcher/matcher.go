package matcher

type Matcher[T any, M any] interface {
	Matches(value *M) bool
	GetValue() *T
}
