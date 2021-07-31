package config

import (
	"testing"
)

type YamlStruct struct {
	Property  string `config:"property"`
	Property2 string `config:"default"`
	User      string `config:"app.db.user"`
	Password  string `config:"app.db.password"`
	Index     int64  `config:"app.db.second[1].data"`
}

type YamlStruct2 struct {
	Property string `config:"property"`
	App      struct {
		Db struct {
			User     string `config:"user"`
			Password string `config:"password"`
		} `config:"db"`
	} `config:"app"`
}

func TestYamlConfigSource(t *testing.T) {

	// d := os.DirFS("tests")
	// csp := &ConfigSourceProvider{}
	// err := csp.AddYaml(d)
	// assert.Nil(t, err)

	// assert.Equal(t, "NO_VALUE", csp.Property("db.user", "NO_VALUE"))
	// assert.Equal(t, "test_user", csp.Property("app.db.user", "NO_VALUE"))
	// assert.Equal(t, "test_password", csp.Property("app.db.password", "NO_VALUE"))
	// assert.Equal(t, "789", csp.Property("app.db.second[1].data", "NO_VALUE"))
	// assert.Equal(t, 789, csp.PropertyInt("app.db.second[1].data", 0))

	// assert.Equal(t, "NO_VALUE", csp.Property("no-property", "NO_VALUE"))
	// assert.Equal(t, "test1", csp.Property("property", "NO_VALUE"))

	// input := &YamlStruct{Property2: "default_value"}
	// err = csp.Properties(input)
	// assert.Nil(t, err)
	// assert.Equal(t, "test1", input.Property)
	// assert.Equal(t, "test_password", input.Password)
	// assert.Equal(t, "test_user", input.User)

	// input2 := &YamlStruct2{}
	// err = csp.Properties(input2)
	// assert.Nil(t, err)
	// assert.Equal(t, "test1", input2.Property)
	// assert.Equal(t, "test_password", input2.App.Db.Password)
	// assert.Equal(t, "test_user", input2.App.Db.User)

	// csp.SetProfile("dev")
	// assert.Equal(t, "test_password", csp.Property("app.db.password", "NO_VALUE"))
	// assert.Equal(t, "test_user_dev", csp.Property("app.db.user", "NO_VALUE"))
	// assert.Equal(t, "test1-dev", csp.Property("property", "NO_VALUE"))

	// err = csp.Properties(input)
	// assert.Nil(t, err)
	// assert.Equal(t, "test_user_dev", input.User)
	// assert.Equal(t, "test_password", input.Password)

	// err = csp.Properties(input2)
	// assert.Nil(t, err)
	// assert.Equal(t, "test1-dev", input2.Property)
	// assert.Equal(t, "test_user_dev", input2.App.Db.User)
	// assert.Equal(t, "test_password", input2.App.Db.Password)

}
