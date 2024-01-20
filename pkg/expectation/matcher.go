package expectation

import "reflect"

type Matcher[T any] struct {
	Match string `json:"match"`
	Value *T     `json:"value"`
}

func (m *Matcher[T]) Matches(value *T) bool {
	switch m.Match {
	case "eq":
		return reflect.DeepEqual(value, m.Value)
	case "has":
		return true
	case "hasValue":
		return value != nil
	}
	// TODO: account for an invalid matcher
	return false
}
