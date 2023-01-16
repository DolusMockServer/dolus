package expectation

type Response struct {
	Body   any
	Status int
}

type Expectation struct {
	Path     string
	Pririoty int
	Response any
}
