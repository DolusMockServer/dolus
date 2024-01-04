package api

// TODO: make this error part of api spec
type GeneralError struct {
	Path     string
	Method   string
	ErrorMsg string
}
