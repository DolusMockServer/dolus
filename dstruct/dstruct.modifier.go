package dstruct

// import (
// 	"fmt"
// 	"reflect"
// 	"strings"

// 	"github.com/MartinSimango/dolus/internal/helper"
// )

// // An extension of "github.com/ompluscator/dynamic-struct" to help modify dynamic structs

// type Field struct {
// 	Name            string
// 	Kind            reflect.Kind // TODO remove
// 	SubFields       []*Field
// 	Tags            reflect.StructTag
// 	Value           reflect.Value
// 	structFieldName string
// }

// type FieldMap map[string]*Field

// type FieldModifier func(*Field)

// type DynamicStructModifierImpl struct {
// 	fields                       FieldMap // holds the meta data for the fields
// 	allFields                    FieldMap
// 	_struct                      any // the struct that actually stores the data
// 	fieldModifier                FieldModifier
// 	dynamicStructModifierBuilder dynamicStructModifierBuilder
// }

// func New(dstruct any) *DynamicStructModifierImpl {
// 	return NewDynamicStructModifierWithFieldModifier(dstruct, nil)
// }

// func NewDynamicStructModifierWithFieldModifier(dstruct any, fieldModifier FieldModifier) *DynamicStructModifierImpl {
// 	ds := &DynamicStructModifierImpl{
// 		_struct:       dstruct,
// 		fields:        make(FieldMap),
// 		allFields:     make(FieldMap),
// 		fieldModifier: fieldModifier,
// 	}
// 	ds.populateFieldMap(ds._struct, "", ds.allFields)

// 	return ds
// }

// // TODO clean this up
// func (ds *DynamicStructModifierImpl) populateFieldMap(config any, root string, allFields FieldMap) (newFields []*Field) {
// 	if config == nil {
// 		return
// 	}

// 	inputConfig := reflect.ValueOf(config).Elem()

// 	for i := 0; i < inputConfig.NumField(); i++ {
// 		field := inputConfig.Field(i)
// 		fieldTags := inputConfig.Type().Field(i).Tag
// 		fieldName := strings.Split(fieldTags.Get("json"), ",")[0]
// 		if fieldName == "" {
// 			fieldName = inputConfig.Type().Field(i).Name
// 		}
// 		if root != "" {
// 			fieldName = fmt.Sprintf("%s.%s", root, fieldName)
// 		}

// 		switch field.Kind() {
// 		case reflect.Struct:
// 			subStruct := &DynamicStructModifierImpl{
// 				_struct:       field.Addr().Interface(),
// 				fields:        make(FieldMap),
// 				fieldModifier: ds.fieldModifier,
// 			}
// 			ds.fields[fieldName] = &Field{
// 				Name:            fieldName,
// 				Kind:            inputConfig.Field(i).Kind(),
// 				SubFields:       subStruct.populateFieldMap(subStruct._struct, fieldName, allFields),
// 				Tags:            fieldTags,
// 				Value:           field,
// 				structFieldName: inputConfig.Type().Field(i).Name,
// 			}
// 		default:
// 			ds.fields[fieldName] = &Field{
// 				Name:            fieldName,
// 				Kind:            inputConfig.Field(i).Kind(),
// 				Tags:            fieldTags,
// 				Value:           field,
// 				structFieldName: inputConfig.Type().Field(i).Name,
// 			}
// 			if ds.fieldModifier != nil {
// 				ds.fieldModifier(ds.fields[fieldName])
// 			}
// 		}

// 		allFields[fieldName] = ds.fields[fieldName]
// 		newFields = append(newFields, ds.fields[fieldName])

// 	}
// 	return

// }

// func (ds *DynamicStructModifierImpl) Get() any {
// 	return helper.GetUnderlyingPointerValue(ds._struct)
// }

// func (ds *DynamicStructModifierImpl) SetFieldValue(field string, value any) error {
// 	f := ds.allFields[field]
// 	if f == nil {
// 		return fmt.Errorf("no such field '%s' exists in schema", field)
// 	}
// 	if value == nil {
// 		f.Value.Set(reflect.Zero(f.Value.Type()))
// 	} else {
// 		f.Value.Set(reflect.ValueOf(value))
// 	}
// 	return nil
// }

// func (ds *DynamicStructModifierImpl) GetFieldValue(field string) (any, error) {
// 	f := ds.allFields[field]
// 	if f == nil {
// 		return nil, fmt.Errorf("no such field '%s' existsx in schema", field)
// 	}

// 	switch f.Kind {
// 	case reflect.String:
// 		return f.Value.String(), nil
// 	case reflect.Int64:
// 		return f.Value.Int(), nil
// 	case reflect.Slice:
// 		sliceLen := reflect.ValueOf(f.Value.Interface()).Len()
// 		return f.Value.Slice(0, sliceLen).Interface(), nil
// 	default:
// 		panic(fmt.Sprintf("unsupported type '%s'", f.Kind))
// 	}
// }

// func (ds *DynamicStructModifierImpl) DoesFieldExist(field string) bool {
// 	return ds.allFields[field] != nil
// }

// func (ds *DynamicStructModifierImpl) GetField(field string) *Field {
// 	return ds.allFields[field]
// }

// func (ds *DynamicStructModifierImpl) DeleteField(field string) {
// 	// dsb := dynamicstruct.ExtendStruct(ds._struct)
// 	// fmt.Println(reflect.TypeOf(dsb.Build().New()))

// 	// fields := strings.Split(field, ".")
// 	// if len(fields) == 1 {
// 	// 	fmt.Println("HERE")
// 	// 	if field == "pair" {
// 	// 		ds.allFields["pair"].SubFields
// 	// 	}
// 	// 	dsb.RemoveField(ds.allFields[fields[0]].structFieldName)
// 	// 	fmt.Println(reflect.TypeOf(dsb.Build().New()))

// 	// 	return
// 	// }

// 	// for _, v := range fields {
// 	// 	ss := dynamicstruct.ExtendStruct(ds.allFields[v].Value.Interface())
// 	// 	fmt.Println(reflect.TypeOf(ss.Build().New()))
// 	// 	continue
// 	// }

// 	// fmt.Println(reflect.TypeOf(dsb.Build().New()))
// 	// fmt.Println()

// }

// func (ds *DynamicStructModifierImpl) Print() {
// 	for k, v := range ds.fields {
// 		fmt.Println("F: ", k, v.Tags, v.Value.Type())
// 	}
// }

// // func setField[T any](f *Field, field string, value T) {
// // 	*(*T)(unsafe.Pointer(f.Address)) = value
// // }

// // func getField[T any](f *Field, field string) T {
// // 	return *(*T)(unsafe.Pointer(f.Address))
// // }
