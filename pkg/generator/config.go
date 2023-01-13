package generator

import "time"

type IntConfig struct {
	Int64Min int64
	Int64Max int64
	Int32Min int32
	Int32Max int32
}

type FloatConfig struct {
	Float64Min float64
	Float64Max float64
	Float32Min float32
	Float32Max float32
}

type SliceConfig struct {
	SliceMinLength int
	SliceMaxLength int
}

type DateConfig struct {
	DateFormat string
	DateStart  time.Time
	DateEnd    time.Time
}

type GenerationConfig struct {
	DefaultGenerationFunctions GenerationDefaults
	ValueGenerationType        ValueGenerationType
	SetNonRequiredFields       bool
	SliceConfig
	IntConfig
	FloatConfig
	DateConfig
}

func defaultIntConfig() IntConfig {
	return IntConfig{
		Int64Min: 0,
		Int64Max: 100,
		Int32Min: 0,
		Int32Max: 100,
	}
}

func defaultFloatConfig() FloatConfig {
	return FloatConfig{
		Float64Min: 0,
		Float64Max: 100,
		Float32Min: 0,
		Float32Max: 100,
	}
}

func defaultSliceConfig() SliceConfig {
	return SliceConfig{
		SliceMinLength: 0,
		SliceMaxLength: 10,
	}
}

func defaultDateConfig() DateConfig {
	return DateConfig{
		DateFormat: time.RFC3339,
	}
}
