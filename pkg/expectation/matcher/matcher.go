package matcher

// TODO: consider merging this package with the expectation package as expectation package cannot use matcher as it causes a circular dependency
type Matcher[T any, M any] interface {
	Matches(value *M) bool
	GetValue() *T
}
