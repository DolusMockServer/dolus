package matcher

type MatcherBuilder[T any] interface {
	Create(field map[string]any) (Matcher[T], error)
	CreateFromArrayValue(value []any, matchExpr string) (Matcher[T], error)
	CreateFromSingleValue(value any, matchExpr string) (Matcher[T], error)
}
