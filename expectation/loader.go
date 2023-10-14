package expectation

type ExpectationType interface {
	OpenAPISpecLoadType | CueExpectationLoadType
}

type Loader[T ExpectationType] interface {
	load() (*T, error)
}
