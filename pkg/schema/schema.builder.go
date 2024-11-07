package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/DolusMockServer/dolus/pkg/task"
)

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

	tags := fmt.Sprintf(
		`json:"%s" type:"%s" pattern:"%s" required:"%s" nullable:"%s" format:"%s"`,
		name,
		schema.Type,
		schema.Pattern,
		required,
		nullable,
		schema.Format,
	)

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

func addField[T any](name string, tags string, nullable bool, builder dstruct.Builder) {
	if nullable {
		builder.AddField(name, new(T), tags)
	} else {
		builder.AddField(name, *new(T), tags)
	}
}

func getStructFromOpenApi3Schema(schema openapi3.Schema) any {
	dsb := dstruct.NewBuilder()

	for name, p := range schema.Properties {
		exportName := getExportName(name)
		tags, nullable := getTags(name, inList(name, schema.Required), *p.Value)
		switch p.Value.Type {
		case "object":
			internalStruct := getStructFromOpenApi3Schema(*p.Value)
			dsb.AddField(exportName, reflect.ValueOf(internalStruct).Interface(), tags)
		case "string":
			addField[string](exportName, tags, nullable, dsb)
		case "number":
			// TODO also check format
			addField[float64](exportName, tags, nullable, dsb)
		case "integer":
			if p.Value.Format == "int32" {
				addField[int32](exportName, tags, nullable, dsb)
			} else {
				addField[int64](exportName, tags, nullable, dsb)
			}
		case "array":
			// TODO need to accomdate
			fmt.Printf("ARRAY: %+v\n", p.Value) //#need to get array items
			panic("ARRAY unimplemented")
		case "boolean":
			addField[bool](exportName, tags, nullable, dsb)
		default:
			panic(fmt.Sprintf("Unsupported type '%s'", p.Value.Type))
		}

	}
	return dsb.Build().Instance()
}

func getStructFromOpenApi3Example(example openapi3.Examples) any {
	for _, v := range example {
		s, _ := buildSchemaFromMap((v.Value.Value).(map[string]interface{}))
		return s
	}
	return nil
}

func getStructFromAny(config any) any {
	schema, _ := buildSchemaFromAny(config, "", "")
	return schema
}

func buildSchemaFromDolusTask(t string, _map map[string]interface{}) (any, SchemaInfo) {
	switch t {
	case task.GenInt32:
		return int32(0), SchemaInfo{
			Kind: reflect.Int32,
			Tags: string(
				generator.GetTask(t).
					GetTags(fmt.Sprintf("%v", _map["min"]), fmt.Sprintf("%v", _map["max"])),
			), // string(generator.GetTagForTask(t, _map["min"], _map["max"])),
			Format: "int32",
		}
	}
	panic(fmt.Sprintf("Unrecognized dolus task: %s", t))
}

// TODO return tags instead of just type within a struct
func buildSchemaFromMap(_map any) (any, SchemaInfo) {
	dsb := dstruct.NewBuilder()
	m := _map.(map[string]interface{})
	for k, v := range m {
		if m["$dolusTask"] != nil {
			task := m["$dolusTask"].(string)
			return buildSchemaFromDolusTask(task, m)
		}

		exportName := getExportName(k)

		schema, schemaInfo := buildSchemaFromAny(v, k, "")
		tags := fmt.Sprintf("%s ", createFieldTags(k, schemaInfo))

		switch schemaInfo.Kind {
		case reflect.String:
			dsb.AddField(exportName, schema.(string), tags)
		case reflect.Float64:
			dsb.AddField(exportName, schema.(float64), tags)
		case reflect.Int:
			dsb.AddField(exportName, schema.(int), tags)
		case reflect.Int32:
			dsb.AddField(exportName, schema.(int32), tags)
		case reflect.Int64:
			dsb.AddField(exportName, int64(schema.(int)), tags)
		case reflect.Bool:
			dsb.AddField(exportName, schema.(bool), tags)
		case reflect.Struct, reflect.Slice:
			dsb.AddField(exportName, schema, tags)
		}
	}
	return dsb.Build().Instance(), SchemaInfo{Kind: reflect.Struct}
}

func buildSchemaFromSlice(config any, name string) (any, SchemaInfo) {
	// fullFieldName := name
	// if root != "" {
	// 	fullFieldName = fmt.Sprintf("%s.%s", root, name)
	// }
	slice := config.([]interface{})

	var firstElement any
	if len(slice) == 0 {
		firstElement = "" // emtpy slice assume array of strings
	} else {
		firstElement, _ = buildSchemaFromAny(slice[0], name, "")
	}

	currentElement := firstElement
	for i := 1; i < len(slice); i++ {
		nextElement, _ := buildSchemaFromAny(slice[i], name, "")
		if reflect.ValueOf(nextElement).Kind() == reflect.Struct {
			var err error

			if currentElement, err = dstruct.MergeStructs(currentElement, nextElement); err != nil {
				panic(err.Error())
			}

			// Account for different types of elements that are not struct
		} else if reflect.TypeOf(nextElement) != reflect.TypeOf(firstElement) {
			// panic(fmt.Sprintf("C: %+v \n\nN:  %+v\n\n", reflect.TypeOf(nextElement), reflect.TypeOf(firstElement)))

			currentElement = ""
			if reflect.ValueOf(firstElement).Kind() == reflect.Slice {
				currentElement = []string{}
			}
			break
		}
	}
	sliceOfElementType := reflect.SliceOf(reflect.ValueOf(currentElement).Type())
	return reflect.MakeSlice(sliceOfElementType, 0, 1024).
			Interface(),
		SchemaInfo{
			Kind: reflect.Slice,
		}
}

func buildSchemaFromStruct(config any, root string) (any, SchemaInfo) {
	dsb := dstruct.NewBuilder()

	inputConfig := reflect.ValueOf(config)
	for i := 0; i < inputConfig.NumField(); i++ {
		field := inputConfig.Field(i)
		fieldTags := inputConfig.Type().Field(i).Tag
		fieldName := strings.Split(fieldTags.Get("json"), ",")[0]

		if fieldName == "" {
			fieldName = inputConfig.Type().Field(i).Name
		}
		if root != "" {
			fieldName = fmt.Sprintf("%s.%s", root, fieldName)
		}

		exportName := getExportName(fieldName)

		schema, schemaInfo := buildSchemaFromAny(field.Interface(), "", "")
		dsb.AddField(exportName, schema, createFieldTags(fieldName, schemaInfo))

	}
	return dsb.Build().Instance(), SchemaInfo{Kind: reflect.Struct}
}

func buildSchemaFromAny(config interface{}, name string, root string) (interface{}, SchemaInfo) {
	if config == nil {
		return nil, SchemaInfo{Kind: reflect.Invalid}
	}
	configKind := reflect.ValueOf(config).Kind()

	switch configKind {
	case reflect.Map:
		return buildSchemaFromMap(config)
	case reflect.Slice:
		return buildSchemaFromSlice(config, name)
	case reflect.Struct:
		return buildSchemaFromStruct(config, root)
	case reflect.Ptr:
		return buildSchemaFromAny(dreflect.GetUnderlyingPointerValue(config), name, root)
	default:
		return config, SchemaInfo{Kind: configKind, Tags: fmt.Sprintf(`default:"%v"`, config)}
	}
}

func getSchemaFromOpenApi3Spec(ref *openapi3.ResponseRef, mediaType string) any {
	content := ref.Value.Content.Get(mediaType)
	if content != nil { // TODO if no example response maybe be empty
		if content.Schema != nil {
			fmt.Println("Using Schema.")
			return getStructFromOpenApi3Schema(*content.Schema.Value)
		} else {
			return getStructFromOpenApi3Example(content.Examples)
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

func createFieldTags(fieldName string, schemaInfo SchemaInfo) string {
	return fmt.Sprintf(`json:"%s" type:"%s" %s`, fieldName, schemaInfo.Kind, schemaInfo.Tags)
}
