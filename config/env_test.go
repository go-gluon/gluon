package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EnvStruct struct {
	Value string `config:"simple.env"`
}

type EnvStruct2 struct {
	Simple struct {
		Env string `config:"env"`
	} `config:"simple"`
}

func TestEnvConfigSource(t *testing.T) {

	os.Setenv("SIMPLE_ENV", "1234")
	os.Setenv("ENV_PROP", "1")
	os.Setenv("_TEST_SIMPLE_ENV", "5678")

	csp := &ConfigSourceProvider{}
	err := csp.Add(&EnvConfigSource{})
	assert.Nil(t, err)

	assert.Equal(t, "1234", csp.Property("SIMPLE_ENV", "NO_VALUE"))
	assert.Equal(t, "1234", csp.Property("simple.env", "NO_VALUE"))
	assert.Equal(t, "1", csp.Property("ENV_PROP", "NO_VALUE"))
	assert.Equal(t, "1", csp.Property("env.prop", "NO_VALUE"))
	assert.Equal(t, "NO_VALUE", csp.Property("SIMPLE_ENV0", "NO_VALUE"))
	assert.Equal(t, "NO_VALUE", csp.Property("simple.env0", "NO_VALUE"))

	s1 := &EnvStruct{}
	err = csp.Properties(s1)
	assert.Nil(t, err)
	assert.Equal(t, "1234", s1.Value)

	s2 := &EnvStruct2{}
	err = csp.Properties(s2)
	assert.Nil(t, err)
	assert.Equal(t, "1234", s2.Simple.Env)

	csp.SetProfile("test")
	assert.Equal(t, "5678", csp.Property("SIMPLE_ENV", "NO_VALUE"))
	assert.Equal(t, "5678", csp.Property("simple.env", "NO_VALUE"))
	assert.Equal(t, "1", csp.Property("ENV_PROP", "NO_VALUE"))
	assert.Equal(t, "1", csp.Property("env.prop", "NO_VALUE"))

	err = csp.Properties(s1)
	assert.Nil(t, err)
	assert.Equal(t, "5678", s1.Value)

	err = csp.Properties(s2)
	assert.Nil(t, err)
	assert.Equal(t, "5678", s2.Simple.Env)
}
