package generator

import (
	"reflect"
)

func setValue(val reflect.Value, tags reflect.StructTag, generationConfig GenerationConfig) {
	switch val.Kind() {
	case reflect.Struct:
		setStructValues(val, generationConfig)
	case reflect.Slice:
		panic("Unhanled setValue case")
	case reflect.Pointer:
		panic("Unhanled setValue case")
	default:
		val.Set(reflect.ValueOf(generationFunctionFromTags(val.Kind(), tags, generationConfig).Generate()))

		// val.SetString(GenerateStringFromRegexFunc("^[a-z ,.'-]+$").Generate().(string))
	}
}

func setStructValues(config reflect.Value, generationConfig GenerationConfig) {
	for j := 0; j < config.NumField(); j++ {
		setValue(config.Field(j), config.Type().Field(j).Tag, generationConfig)
	}

}
