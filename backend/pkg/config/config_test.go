package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestNew_WithDefaultConfigFile(t *testing.T) {
	Init(New(WithDefaultConfigFile("", "")))
	tName := GetString("worker.inputs.name")
	assert.Equal(t, tName, "ABC")
}

func TestUnmarshalKey_FromEnv(t *testing.T) {
	os.Setenv("WORKER_INPUT", "{\"name\":\"ABC\",\"age\":11}")
	os.Setenv("WORKER_INPUTS", "[{\"name\":\"ABC\",\"age\":11}]")
	os.Setenv("WORKER_INPUTM", "{\"ABC\":{\"name\":\"ABC\",\"age\":11}}")

	Init(New(WithDefaultEnvVars("")))

	var tss []*testStruct
	UnmarshalKey("worker.inputs", &tss)
	require.Equal(t, 1, len(tss))
	assert.Equal(t, tss[0].Name, "ABC")
	assert.Equal(t, tss[0].Age, 11)

	var ts = &testStruct{}
	UnmarshalKey("worker.input", &ts)
	assert.Equal(t, ts.Name, "ABC")
	assert.Equal(t, ts.Age, 11)

	var tms = map[string]testStruct{}
	UnmarshalKey("worker.inputm", &tms)
	require.NotNil(t, tms["ABC"])
	assert.Equal(t, tms["ABC"].Name, "ABC")
	assert.Equal(t, tms["ABC"].Age, 11)
}

func TestUnmarshalKey_FromYaml(t *testing.T) {
	yamlContent := []byte(`
worker:
  inputs:
    name: ABC
    age: 11
`)
	Init(New(WithReader("yaml", bytes.NewBuffer(yamlContent))))

	ts := &testStruct{}
	UnmarshalKey("worker.inputs", ts)
	assert.Equal(t, ts.Name, "ABC")
	assert.Equal(t, ts.Age, 11)
}

func TestGetFloat32Slice_FromEnv(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "0 1.2E-5 3.2 4.1")
	Init(New(WithDefaultEnvVars("")))
	a := GetFloat32Slice("sample.input")
	assert.Equal(t, []float32{0, 1.2e-05, 3.2, 4.1}, a)
}

func TestGetFloat32Slice_FromYaml(t *testing.T) {
	yamlContent := []byte(`a: ["0", 1.2E-5, 3.2, 4.1]`)
	Init(New(WithReader("yaml", bytes.NewBuffer(yamlContent))))
	a := GetFloat32Slice("a")
	assert.Equal(t, []float32{0, 1.2e-05, 3.2, 4.1}, a)
}

func TestGetFloat32Slice_Error(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "\"abc\" 1.2E-5 3.2 4.1")
	Init(New(WithDefaultEnvVars("")))
	a := GetFloat32Slice("sample.input")
	assert.Equal(t, []float32{}, a)
}

func TestGetFloat64Slice_FromEnv(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "0 1.2E-5 3.2 4.1")
	Init(New(WithDefaultEnvVars("")))
	a := GetFloat64Slice("sample.input")
	assert.Equal(t, []float64{0, 1.2e-05, 3.2, 4.1}, a)
}

func TestGetFloat64Slice_FromYaml(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "0 1.2E-5 3.2 4.1")
	Init(New(WithDefaultEnvVars("")))
	a := GetFloat64Slice("sample.input")
	assert.Equal(t, []float64{0, 1.2e-05, 3.2, 4.1}, a)
}

func TestGetFloat64Slice_Error(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "\"abc\" 1.2E-5 3.2 4.1")
	Init(New(WithDefaultEnvVars("")))
	a := GetFloat64Slice("sample.input")
	assert.Equal(t, []float64{}, a)
}

func TestGetIntSlice_FromEnv(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "1 2 3 4")
	Init(New(WithDefaultEnvVars("")))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{1, 2, 3, 4}, intSlice)
}

func TestGetIntSlice_FromYaml(t *testing.T) {
	yamlContent := []byte(`
sample:
  input:
    - 1
    - 2
    - 3
    - 4
`)
	Init(New(WithReader("yaml", bytes.NewBuffer(yamlContent))))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{1, 2, 3, 4}, intSlice)
}

func TestGetInt64Slice_Error(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "1 fail 3 4")
	Init(New(WithDefaultEnvVars("")))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{}, intSlice)
}

func TestGetInt64Slice_FromEnv(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "1 2 3 4")
	Init(New(WithDefaultEnvVars("")))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{1, 2, 3, 4}, intSlice)
}

func TestGetInt64Slice_FromYaml(t *testing.T) {
	yamlContent := []byte(`
sample:
  input:
    - 1
    - 2
    - 3
    - 4
`)
	Init(New(WithReader("yaml", bytes.NewBuffer(yamlContent))))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{1, 2, 3, 4}, intSlice)
}

func TestGetIntSlice_Error(t *testing.T) {
	os.Setenv("SAMPLE_INPUT", "1 fail 3 4")
	Init(New(WithDefaultEnvVars("")))
	intSlice := GetIntSlice("sample.input")
	assert.Equal(t, []int{}, intSlice)
}
