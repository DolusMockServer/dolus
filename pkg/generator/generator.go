package generator

import (
	"reflect"
)

type GenerationFunction interface {
	Generate() any
}

type GenerationFunc struct {
	_func func(...any) any
	args  []any
}

var _ GenerationFunction = &GenerationFunc{}

func (f GenerationFunc) Generate() any {
	return f._func(f.args...)
}

type ValueGenerationType uint8

const (
	Generate     ValueGenerationType = iota // will generate all field
	GenerateOnce                            // will generate all the fields
	UseDefaults
)

type GenerationDefaults map[reflect.Kind]GenerationFunction

func NewGenerationConfig() (generationConfig *GenerationConfig) {
	generationConfig = &GenerationConfig{
		ValueGenerationType:  Generate,
		SetNonRequiredFields: false,
		SliceConfig:          defaultSliceConfig(),
		IntConfig:            defaultIntConfig(),
		FloatConfig:          defaultFloatConfig(),
		DateConfig:           defaultDateConfig(),
	}
	generationConfig.initGenerationFunctionDefaults()
	return
}

func NewGenerationConfigFromConfig(generationConfig GenerationConfig) (config *GenerationConfig) {
	config = &GenerationConfig{}
	*config = generationConfig
	config.initGenerationFunctionDefaults()
	return config
}

func (gc *GenerationConfig) initGenerationFunctionDefaults() {
	gc.DefaultGenerationFunctions = make(GenerationDefaults)
	gc.DefaultGenerationFunctions[reflect.String] = GenerateFixedValueFunc("string") //generator.GenerateStringFromRegexFunc("^[a-z ,.'-]+$")
	gc.DefaultGenerationFunctions[reflect.Ptr] = GenerateNilValue()
	gc.DefaultGenerationFunctions[reflect.Int64] = GenerateInt64Func(&gc.Int64Min, &gc.Int64Max)
	gc.DefaultGenerationFunctions[reflect.Int32] = GenerateInt32Func(&gc.Int32Min, &gc.Int32Max)
	gc.DefaultGenerationFunctions[reflect.Float64] = GenerateFloatFunc(&gc.Float64Min, &gc.Float64Max)
	gc.DefaultGenerationFunctions[reflect.Bool] = GenerateBoolFunc()

}

type GenerationUnit struct {
	CurrentFunction  GenerationFunction
	GenerationConfig *GenerationConfig
	count            int
	latestValue      any
}

func (g *GenerationUnit) Generate() any {
	if g.GenerationConfig.ValueGenerationType == GenerateOnce && g.count > 0 {
		return g.latestValue
	}
	g.latestValue = g.CurrentFunction.Generate()
	g.count++
	return g.latestValue
}
