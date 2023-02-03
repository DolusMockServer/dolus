package expectation

//cue:generate cue get go
type CueResponse struct {
	Body   any `json:"body"`
	Status int `json:"status"`
}

type CueRequest struct {
	Body any `json:"body"`
}

type CueExpectation struct {
	Path     string      `json:"path"`
	Method   string      `json:"method"`
	Pririoty int         `json:"priority"`
	Response CueResponse `json:"response"`
	Request  CueRequest  `json:"request"`
}
