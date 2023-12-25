package loader

type ExpectationType interface {
	OpenAPISpecLoadType | CueExpectationLoadType
}

type Loader[T ExpectationType] interface {
	Load() (*T, error)
}
