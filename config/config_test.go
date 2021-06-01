package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigStruct struct {
	Property  string `config:"property"`
	Property2 string `config:"default"`
	User      string `config:"app.db.user"`
	Password  string `config:"app.db.password"`
	Index     int64  `config:"app.db.second[1].data"`
}

type ConfigStruct2 struct {
	Property string `config:"property"`
	App      struct {
		Db struct {
			User     string `config:"user"`
			Password string `config:"password"`
		} `config:"db"`
	} `config:"app"`
}

type TestExtensionConfig struct {
	Config interface{}
}

type Config struct {
	Name string `config:"name"`
}

func TestExtension(t *testing.T) {

	tec := TestExtensionConfig{
		Config: &Config{},
	}

	d := os.DirFS(".")
	err := Default.AddYaml(d)
	assert.Nil(t, err)

	err = Default.Extension("test", tec.Config)
	assert.Nil(t, err)
}

func TestInput(t *testing.T) {

	d := os.DirFS(".")
	err := Default.AddYaml(d)
	assert.Nil(t, err)

	assert.Equal(t, "NO_VALUE", Default.Property("db.user", "NO_VALUE"))
	assert.Equal(t, "test_user", Default.Property("app.db.user", "NO_VALUE"))
	assert.Equal(t, "789", Default.Property("app.db.second[1].data", "NO_VALUE"))
	assert.Equal(t, 789, Default.PropertyInt("app.db.second[1].data", 0))

	assert.Equal(t, "NO_VALUE", Default.Property("no-property", "NO_VALUE"))
	assert.Equal(t, "test1", Default.Property("property", "NO_VALUE"))

	input := &ConfigStruct{Property2: "default_value"}
	err = Default.Properties(input)
	assert.Nil(t, err)
	t.Logf("Input %v\n", input)

	input2 := &ConfigStruct2{}
	err = Default.Properties(input2)
	assert.Nil(t, err)
	t.Logf("Input2 %v\n", input2)
}
