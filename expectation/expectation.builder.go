package expectation

type ExpectationBuilder interface {
	BuildExpectationsFromCueLoader(loader Loader[CueExpectationLoadType]) ([]Expectation, error)
	BuildExpectationsFromOpenApiSpecLoader(loader Loader[OpenAPISpecLoadType]) ([]Expectation, error)
}
