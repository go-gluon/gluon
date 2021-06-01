package config

import (
	"errors"
	"io/fs"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

const resourceFile = "application.yaml"

// AddYaml adds yaml configuration source for the embedded yaml file to the default provider
func AddYaml(resources fs.FS) error {
	return Default.AddYaml(resources)
}

// AddYaml adds yaml configuration source for the embedded yaml file
func (c *ConfigSourceProvider) AddYaml(resources fs.FS) error {
	return c.Add(&YamlConfigSource{resources: resources, Prio: 100})
}

type YamlConfigSource struct {
	resources fs.FS
	Prio      int
	data      map[string]string
}

var errorFind = errors.New("find item")

func (y *YamlConfigSource) Init() error {
	y.data = map[string]string{}

	rf := resourceFile
	er := fs.WalkDir(y.resources, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if d.Name() == resourceFile {
				rf = path
				return errorFind
			}
		}
		return nil
	})
	if er != nil && er != errorFind {
		return er
	}

	d, err := fs.ReadFile(y.resources, rf)
	if err != nil {
		return err
	}

	tmp := map[string]interface{}{}
	err = yaml.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	flatten(tmp, "", y.data)
	return nil
}

func (y *YamlConfigSource) Priority() int {
	return y.Prio
}

func (y *YamlConfigSource) Name() string {
	return `yaml`
}

func (y *YamlConfigSource) Property(name string) (string, bool, error) {
	v, e := y.data[name]
	return v, e, nil
}

func (y *YamlConfigSource) Properties() (map[string]string, error) {
	return y.data, nil
}

func flatten(value interface{}, prefix string, m map[string]string) {

	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}
	t := original.Type()

	switch kind {
	case reflect.Map:
		for _, childKey := range original.MapKeys() {
			flatten(original.MapIndex(childKey).Interface(), flattenName(prefix, toString(childKey)), m)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < original.Len(); i++ {
			flatten(original.Index(i).Interface(), flattenArray(prefix, i), m)
		}
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			flatten(original.Field(i).Interface(), flattenName(prefix, t.Field(i).Name), m)
		}
	default:
		if prefix != "" {
			m[prefix] = toString(value)
		}
	}
}

func flattenArray(prefix string, index int) string {
	return prefix + "[" + strconv.Itoa(index) + "]"
}

func flattenName(prefix, child string) string {
	if prefix != "" {
		return prefix + "." + child
	}
	return child
}
