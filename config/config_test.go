package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ConfigStruct struct {
	Property  string `config:"property"`
	Property2 string `config:"default"`
	User      string `config:"app.db.user"`
	Password  string `config:"app.db.password"`
	Index     int64  `config:"app.db.second[1].data"`
}

type ConfigStruct2Data struct {
	Name  string `config:"name"`
	Value int    `config:"value"`
}
type ConfigStruct2 struct {
	Property string `config:"property"`
	App      struct {
		Db struct {
			User     string                       `config:"user"`
			Password string                       `config:"password"`
			List     []string                     `config:"list"`
			Data     map[string]ConfigStruct2Data `config:"data"`
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

	d := os.DirFS(".")
	err := Default.Init(d)
	assert.Nil(t, err)

	c := ConfigStruct2{}
	db := Default.Map().Map("app").Map("db")

	c.App.Db.User = db.String("user", "defaultValue1")
	assert.Equal(t, "test_user", c.App.Db.User)
	c.App.Db.Password = db.String("password", "defaultValue1")
	assert.Equal(t, "test_password", c.App.Db.Password)

	c.App.Db.List = db.StringL("list", nil)
	assert.Equal(t, 2, len(c.App.Db.List))

	data := db.Map("data")
	if c.App.Db.Data == nil {
		c.App.Db.Data = map[string]ConfigStruct2Data{}
	}
	for _, key := range data.Keys() {
		cc := ConfigStruct2Data{}
		item := data.Map(key)
		cc.Name = item.String("name", "")
		cc.Value = item.Int("value", -1)
		c.App.Db.Data[key] = cc
	}
	assert.Equal(t, 2, len(c.App.Db.Data))
	assert.Equal(t, "n1", c.App.Db.Data["key1"].Name)
	assert.Equal(t, 123, c.App.Db.Data["key1"].Value)

	r := Default.Map()
	assert.Equal(t, -9223372036854775808, r.Int("int", 0))
	assert.Equal(t, 1.34, r.Float64("float", 0))
	assert.Equal(t, "2018-01-09 10:40:47 +0000 UTC", fmt.Sprintf("%v", r.Time("time", time.Now())))

	ms := Default.Map().StringM("mapString", map[string]string{})
	assert.Equal(t, 2, len(ms))
	assert.Equal(t, "1", ms["A"])
	assert.Equal(t, "2", ms["B"])
}

func TestConfigKeys(t *testing.T) {

	d := os.DirFS(".")
	err := Default.Init(d)
	assert.Nil(t, err)
	assert.Equal(t, "test_user", Default.String("app.db.user", "defaultValue1"))
	assert.Equal(t, "test_password", Default.String("app.db.password", "defaultValue1"))
	assert.Equal(t, 1.34, Default.Float("float", 0))
	assert.Equal(t, -9223372036854775808, Default.Int("int", 0))
	assert.Equal(t, "2018-01-09 10:40:47 +0000 UTC", fmt.Sprintf("%v", Default.Time("time", time.Now())))

	assert.Equal(t, "n1", Default.String("app.db.data.key1.name", "defaultValue1"))
}

func TestExtensionProfileDev(t *testing.T) {

	d := os.DirFS(".")
	Default.SetProfile("dev")
	err := Default.Init(d)
	assert.Nil(t, err)

	c := ConfigStruct2{}
	db := Default.Map().Map("app").Map("db")
	c.App.Db.User = db.String("user", "defaultValue1")
	assert.Equal(t, "test_user_dev", c.App.Db.User)
}

func TestExtensionProfileTest(t *testing.T) {

	d := os.DirFS(".")
	Default.SetProfile("test")
	err := Default.Init(d)
	assert.Nil(t, err)

	c := ConfigStruct2{}
	db := Default.Map().Map("app").Map("db")
	c.App.Db.User = db.String("user", "defaultValue1")
	assert.Equal(t, "defaultValue1", c.App.Db.User)
}
