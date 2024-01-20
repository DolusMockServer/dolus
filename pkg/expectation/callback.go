package expectation

type Callback struct {
	Request any    `json:"request"`
	Timeout int    `json:"timeout"`
	Url     string `json:"url"`
	Method  string `json:"method"`
}
