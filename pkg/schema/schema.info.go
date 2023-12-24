package schema

import (
	"reflect"
)

type SchemaInfo struct {
	Kind   reflect.Kind
	Tags   string
	Format string
}
