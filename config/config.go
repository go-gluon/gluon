package config

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// ConfigSource configuration source interface
type ConfigSource interface {
	// Init initialize method
	Init() error
	// Name of the configuration source
	Name() string
	// Priority of the configuration source
	Priority() int
	// Properties is the map of all avaiable properties in the configuration source
	Properties() (map[string]string, error)
	// Property get property by the name
	Property(name string) (string, bool, error)
}

const (
	configPrefix = "gluon."
)

var (
	// Default configuration source provider instance
	Default *ConfigSourceProvider
	// ConfigProfileProperty configuration profile property
	ConfigProfileProperty = "config.profile"
)

// init initialize default configuration
func init() {
	// create default configuration source provider
	Default = &ConfigSourceProvider{}

	// environment configuration
	e := &EnvConfigSource{}

	// flags parameter configurations
	f := &FlagsConfigSource{}

	// add default configuration sources
	err := Default.Add(f, e)
	if err != nil {
		panic(err)
	}

	// set profile
	profile := Default.Property(ConfigProfileProperty, "")
	if len(profile) > 0 {
		Default.SetProfile(profile)
	}
}

// ConfigSourceProvider configuration source provider
type ConfigSourceProvider struct {
	sources    []ConfigSource
	profile    string
	profileOrg string
}

// SetProfile set profile to default provider
func SetProfile(profile string) {
	Default.SetProfile(profile)
}

// Profile of the default configuration source provider
func Profile() string {
	return Default.profile
}

// SetProfile set configuration profile to provider
func (c *ConfigSourceProvider) SetProfile(profile string) {
	c.profileOrg = profile
	if len(profile) > 0 {
		c.profile = "+" + profile + "."
	} else {
		c.profile = ""
	}
}

// Profile configuration profile of the configuration source provider
func (c *ConfigSourceProvider) Profile() string {
	return c.profile
}

// Properties setup the properties in the structure base on the tags
func Properties(value interface{}) error {
	return Default.Properties(value)
}

// Extension setup the properties in the structure base on the tags
func Extension(name string, value interface{}) error {
	return Default.Extension(name, value)
}

// Properties setup the properties in the structure base on the tags
func (c *ConfigSourceProvider) Properties(value interface{}) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("Configuration properties is not pointer to struct")
	}
	c.properties("", value)
	return nil
}

// Extension setup the properties in the structure base on the tags
func (c *ConfigSourceProvider) Extension(name string, value interface{}) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("Extension configuration is not pointer to struct")
	}
	c.properties(configPrefix+name, value)
	return nil
}

// add properties to the struct
func (c *ConfigSourceProvider) properties(prefix string, value interface{}) {
	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
	} else {
		return
	}

	typeof := reflect.TypeOf(value)
	if typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
	}

	// fmt.Printf("NumField: %v\n", typeof.NumField())
	for i := 0; i < typeof.NumField(); i++ {
		f := typeof.Field(i)
		prop := f.Name
		tag, ok := f.Tag.Lookup("config")
		if ok {
			prop = tag
		}
		if prefix != "" {
			prop = prefix + "." + tag
		}
		field := original.Field(i)
		switch field.Kind() {
		case reflect.String:
			tmp := c.Property(prop, field.String())
			field.SetString(tmp)
		case reflect.Float32, reflect.Float64:
			tmp := c.Property(prop, toString(field.Float()))
			f, _ := strconv.ParseFloat(tmp, 64)
			field.SetFloat(f)
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			tmp := c.Property(prop, toString(field.Int()))
			f, _ := strconv.ParseInt(tmp, 10, 0)
			field.SetInt(f)
		case reflect.Bool:
			tmp := c.Property(prop, toString(field.Bool()))
			b, _ := strconv.ParseBool(tmp)
			field.SetBool(b)
		case reflect.Struct:

			fval := reflect.New(f.Type)
			c.properties(prop, fval.Interface())
			field.Set(fval.Elem())
		default:
			fmt.Printf("Not supported field: %v type: %v\n", f.Name, field.Kind())
		}

	}

}

// PropertyBool bool value property from the default configuration source provider
func PropertyBool(name string, defaultValue bool) bool {
	return Default.PropertyBool(name, defaultValue)
}

// PropertyBool bool value property from the configuration source provider
func (c *ConfigSourceProvider) PropertyBool(name string, defaultValue bool) bool {
	value, exists := c.findProperty(name)
	if exists {
		tmp, err := strconv.ParseBool(value)
		if err != nil {
			//TODO: debug log
			return defaultValue
		}
		return tmp
	}
	return defaultValue
}

// PropertyInt int value property from the default configuration source provider
func PropertyInt(name string, defaultValue int) int {
	return Default.PropertyInt(name, defaultValue)
}

// PropertyInt int value property from the configuration source provider
func (c *ConfigSourceProvider) PropertyInt(name string, defaultValue int) int {
	value, exists := c.findProperty(name)
	if exists {
		tmp, err := strconv.Atoi(value)
		if err != nil {
			//TODO: debug log
			return defaultValue
		}
		return tmp
	}
	return defaultValue
}

// Property string value property from the default configuration source provider
func Property(name string, defaultValue string) string {
	return Default.Property(name, defaultValue)
}

// Property string value property from the configuration source provider
func (c *ConfigSourceProvider) Property(name, defaultValue string) string {
	value, exists := c.findProperty(name)
	if exists {
		return value
	}
	return defaultValue
}

// Add methods add configuration sources to the default provider
func Add(s ...ConfigSource) error {
	return Default.Add(s...)
}

// Add methods add configuration sources
func (c *ConfigSourceProvider) Add(s ...ConfigSource) error {
	if len(s) > 0 {
		for _, item := range s {
			err := item.Init()
			if err != nil {
				return err
			}
			c.sources = append(c.sources, item)
		}
	}
	sort.SliceStable(c.sources, func(i, j int) bool {
		return c.sources[i].Priority() > c.sources[j].Priority()
	})
	return nil
}

func (c *ConfigSourceProvider) findProperty(name string) (string, bool) {
	if len(c.sources) > 0 {
		for _, source := range c.sources {
			if len(c.profile) > 0 {
				value, exists, err := source.Property(c.profile + name)
				if err != nil {
					//TODO: debug log
					return "", false
				}
				if exists {
					return value, true
				}
			}
			value, exists, err := source.Property(name)
			if err != nil {
				//TODO: debug log
				return "", false
			}
			if exists {
				return value, true
			}
		}
	}
	return "", false
}

func toString(data interface{}) string {
	return fmt.Sprintf("%v", data)
}
