package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"github.com/MartinSimango/dolus/pkg/dstruct"
	"github.com/MartinSimango/dolus/pkg/helper"
	"github.com/getkin/kin-openapi/openapi3"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

type ResponseSchema struct {
	Path       string
	Method     string
	StatusCode string
	schema     any
}

func New(path, method, statusCode string, ref *openapi3.ResponseRef, mediaType string,
) *ResponseSchema {
	return &ResponseSchema{
		Path:       path,
		Method:     method,
		StatusCode: statusCode,
		schema:     getSchema(ref, mediaType),
	}

}
func NewSchemaFromCueValue(path, method, statusCode string, s any) *ResponseSchema {
	return &ResponseSchema{
		Path:       path,
		Method:     method,
		StatusCode: statusCode,
		schema:     s,
	}
}
func (rs *ResponseSchema) GetSchema() any {
	// Make copy of schema to use as struct that is being modified to not modify original schema
	if rs.schema == nil {
		return nil
	}
	schemaValue := reflect.ValueOf(rs.schema).Elem().Interface()
	return reflect.New(reflect.ValueOf(schemaValue).Type()).Interface()
}

// func
func (schema *ResponseSchema) MarshalSchema() (string, error) {
	bytes, err := json.Marshal(schema.schema)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func inList(key string, list []string) bool {
	for k := range list {
		if list[k] == key {
			return true
		}
	}
	return false
}

func getTags(name string, requiredField bool, schema openapi3.Schema) (string, bool) {
	nullable := "true"
	if !schema.Nullable {
		nullable = "false"
	}
	required := "false"
	if requiredField {
		required = "true"
	} else {
		name = fmt.Sprintf("%s,%s", name, "omitempty")
	}
	pointer := nullable == "true" || required == "false"

	tags := fmt.Sprintf(`json:"%s" type:"%s" pattern:"%s" required:"%s" nullable:"%s" format:"%s"`, name,
		schema.Type,
		schema.Pattern,
		required,
		nullable,
		schema.Format)

	if schema.Type != "object" {
		if schema.Example != nil {
			tags = fmt.Sprintf(`%s example:"%v"`, tags, schema.Example)

		}
		if schema.Default != nil {
			tags = fmt.Sprintf(`%s default:"%v"`, tags, schema.Default)
		}
		if schema.Enum != nil {
			len := len(schema.Enum)
			tags = fmt.Sprintf(`%s enum:"%v"`, tags, len)
			for i := 0; i < len; i++ {
				tags = fmt.Sprintf(`%s enum_%d:"%v"`, tags, (i + 1), schema.Enum[i])
			}

		}
	}

	return tags, pointer
}

func addField[T any](name string, tags string, nullable bool, builder *dynamicstruct.Builder) {
	if nullable {
		(*builder).AddField(name, new(T), tags)
	} else {
		(*builder).AddField(name, *new(T), tags)
	}
}

func structFromSchema(schema openapi3.Schema) any {
	dsb := dynamicstruct.NewStruct()

	for name, p := range schema.Properties {
		exportName := helper.GetExportName(name)
		tags, nullable := getTags(name, inList(name, schema.Required), *p.Value)
		switch p.Value.Type {
		case "object":
			internalStruct := structFromSchema(*p.Value)
			dsb.AddField(exportName, reflect.ValueOf(internalStruct).Elem().Interface(), tags)
		case "string":
			addField[string](exportName, tags, nullable, &dsb)
		case "number":
			// TODO also check format
			addField[float64](exportName, tags, nullable, &dsb)
		case "integer":
			if p.Value.Format == "int32" {
				addField[int32](exportName, tags, nullable, &dsb)

			} else {
				addField[int64](exportName, tags, nullable, &dsb)
			}
		case "array":
			// TODO need to accomdate
			fmt.Printf("ARRAY: %+v\n", p.Value) //#need to get array items
		case "boolean":
			addField[bool](exportName, tags, nullable, &dsb)
		default:
			panic(fmt.Sprintf("Unsupported type '%s'", p.Value.Type))
		}

	}
	return dsb.Build().New()
}

func structFromExample(example openapi3.Examples) any {
	// TODO make this not use a loop
	for _, v := range example {
		m := (v.Value.Value).(map[string]interface{})
		g, _ := BuildExample(m, "", "", nil)
		return reflect.New(reflect.ValueOf(g).Type()).Interface()
	}
	fmt.Println()
	return nil
}

func getTagsFromDolusTask(task string, _map map[string]interface{}) (any, Tag) {
	switch task {
	case "GenInt32":
		return float64(0), Tag{Type: "integer",
			Tags:   fmt.Sprintf(`gen_task:"%s" gen_param_1:"%d" gen_param_2:"%d"`, task, int32(_map["min"].(float64)), int32(_map["max"].(float64))),
			Format: "int32"}
	}
	panic(fmt.Sprintf("Unrecognised dolus task: %s", task))
}

type Tag struct {
	Type   string
	Tags   string
	Format string
}

// TODO return tags instead of just type within a struct
func buildStructFromMap(_map any, cueValue *cue.Value) (any, Tag) {
	dsb := dynamicstruct.NewStruct()
	m := _map.(map[string]interface{})
	for k, v := range m {
		var nextCueValue *cue.Value
		if cueValue != nil {
			c := cueValue.Lookup(k)
			nextCueValue = &c
		}
		if m["$dolus"] != nil {
			task := m["$dolus"].(map[string]interface{})["task"].(string)
			return getTagsFromDolusTask(task, m)
		}

		exportName := getExportName(k)

		i, _type := BuildExample(v, k, "", nextCueValue)
		if cueValue != nil {
			switch cueValue.Lookup(k).Kind() {
			case cue.FloatKind:
				_type.Type = "number"
			}
		}
		tags := fmt.Sprintf(`json:"%s" type:"%s" %s`, k, _type.Type, _type.Tags)

		switch _type.Type {
		case "string":
			dsb.AddField(exportName, i.(string), tags)
		case "number":
			dsb.AddField(exportName, i.(float64), tags)
		case "integer":
			if _type.Format == "int32" {
				dsb.AddField(exportName, int32(i.(float64)), tags)
			} else {
				dsb.AddField(exportName, int64(i.(float64)), tags)
			}

		case "boolean":
			dsb.AddField(exportName, i.(bool), tags)
		case "slice":
			dsb.AddField(exportName, i, tags)

		case "struct":
			dsb.AddField(exportName, i, tags)
		}
	}
	return reflect.ValueOf(dsb.Build().New()).Elem().Interface(), Tag{Type: "struct"}
}

func buildSliceOfSliceElementType(config any, name string, root string, cueValue *cue.Value) (any, Tag) {
	fullFieldName := name
	if root != "" {
		fullFieldName = fmt.Sprintf("%s.%s", root, name)
	}
	slice := config.([]interface{})

	var firstElement any
	if len(slice) == 0 {
		firstElement = "" //emtpy slice assume array of strings
	} else {
		firstElement, _ = BuildExample(slice[0], name, "", cueValue)
	}

	currentElement := firstElement
	for i := 1; i < len(slice); i++ {
		nextElement, _ := BuildExample(slice[i], name, "", cueValue)
		if reflect.ValueOf(nextElement).Kind() == reflect.Struct {
			var err error
			var mergedStruct *dstruct.DynamicStructModifier
			if mergedStruct, err = dstruct.MergeStructs(currentElement, nextElement, fullFieldName); err != nil {
				panic(err.Error())
			}
			currentElement = mergedStruct.Get()
			// Account for different types of elements that are not struct
		} else if reflect.TypeOf(nextElement) != reflect.TypeOf(firstElement) {
			currentElement = ""
			if reflect.ValueOf(firstElement).Kind() == reflect.Slice {
				currentElement = []string{}
			}
			break
		}
	}
	sliceOfElementType := reflect.SliceOf(reflect.ValueOf(currentElement).Type())
	return reflect.MakeSlice(sliceOfElementType, 0, 1024).Interface(), Tag{Type: "slice"}
}

func getType(element any, kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Float64:
		val := element.(float64)
		if val-float64(int64(val)) == 0 {
			return "integer"
		} else {
			return "number"
		}

	case reflect.Bool:
		return "boolean"
	}
	return "unknown"
}

func BuildExample(config interface{}, name string, root string, cueValue *cue.Value) (interface{}, Tag) {

	if config == nil {
		return nil, Tag{Type: "nil"}
	}
	configKind := reflect.ValueOf(config).Kind()
	switch configKind {
	case reflect.Map:
		return buildStructFromMap(config, cueValue)
	case reflect.Slice:
		return buildSliceOfSliceElementType(config, name, root, cueValue)
	default:
		return config, Tag{Type: getType(config, configKind)}
	}

}

func getSchema(ref *openapi3.ResponseRef, mediaType string) any {
	content := ref.Value.Content.Get(mediaType)
	if content != nil { // TODO if no example response maybe be empty
		if content.Schema != nil {
			fmt.Println("Using Schema.")
			return structFromSchema(*content.Schema.Value)
		} else {
			return structFromExample(content.Examples)
		}

	}

	fmt.Println("NO SCHEMA!")
	return nil
}

func getExportName(name string) string {
	// TODO replace special characters in name
	name = strings.ReplaceAll(name, "-", "_")
	return strings.ToUpper(name[:1]) + name[1:]
}
