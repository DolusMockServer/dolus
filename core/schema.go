package core

import (
	"encoding/json"
)

type schema interface {
	GetPath() string
	GetMethod() string
	GetSchema() any
	MarshalSchema() (string, error)
}

type Schema struct {
	Path   string
	Method string
	schema any
}

func (s *Schema) GetPath() string {
	return s.Path
}

func (s *Schema) GetMethod() string {
	return s.Method
}

func (s *Schema) GetSchema() any {
	return s.schema
}

func (s *Schema) MarshalSchema() (string, error) {
	bytes, err := json.Marshal(s.schema)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
