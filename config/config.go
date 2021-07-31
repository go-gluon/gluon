package config

import (
	"io/fs"
	"time"
)

const (
	configPrefix = "gluon."
)

var (
	// Default configuration source provider instance
	Default *ConfigSourceProvider
	// ConfigProfileProperty configuration profile property
	configProfileProperty = configPrefix + "config.profile"
)

type ConfigReader interface {
	ReadFromSource(node MapNode)
}

// ConfigSource configuration source interface
type ConfigSource interface {
	// Init initialize method
	Init() error

	GetRawValue(name string) (string, bool, error)
}

// init initialize default configuration
func init() {

	// initialize config sources
	sources := []ConfigSource{&EnvConfigSource{}, &FlagsConfigSource{}}
	for _, s := range sources {
		err := s.Init()
		if err != nil {
			panic(err)
		}
	}

	// create default configuration source provider
	Default = &ConfigSourceProvider{
		data:    map[interface{}]interface{}{},
		sources: sources,
	}

	// set profile
	profile := Default.getRawValue(configProfileProperty, "")
	if len(profile) > 0 {
		Default.SetProfile(profile)
	}
}

// ConfigSourceProvider configuration source provider
type ConfigSourceProvider struct {
	sources []ConfigSource
	data    map[interface{}]interface{}
	profile string
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
	if len(profile) > 0 {
		c.profile = `+` + profile
	} else {
		c.profile = ""
	}
}

// Profile configuration profile of the configuration source provider
func (c *ConfigSourceProvider) Profile() string {
	return c.profile
}

func GetRawValue(name string, default_value string) string {
	return Default.getRawValue(name, default_value)
}

func (c *ConfigSourceProvider) getRawValue(name string, default_value string) string {
	for _, s := range c.sources {
		v, e, _ := s.GetRawValue(name)
		if e {
			return v
		}
	}
	return default_value
}

func (c *ConfigSourceProvider) LoadYaml(resources fs.FS) error {
	data, err := loadYaml(resources)
	if err != nil {
		return err
	}
	c.data = data
	if len(c.profile) > 0 {
		c.mergeProfile()
	}
	return nil
}

func (c *ConfigSourceProvider) mergeProfile() {
	if len(c.profile) == 0 {
		return
	}
	d, e := c.data[c.profile]
	if !e {
		return
	}
	data, ee := d.(map[interface{}]interface{})
	if !ee {
		return
	}
	merge(data, c.data)
}

func merge(profile, config map[interface{}]interface{}) {

	for key, value := range profile {

		// // skip nil values
		// if value == nil {
		// 	continue
		// }

		// config does not exists
		cv, ce := config[key]
		if !ce {
			config[key] = value
			continue
		}

		// value is not map
		v, ok := value.(map[interface{}]interface{})
		if !ok {
			config[key] = value
			continue
		}

		// merge children
		merge(v, cv.(map[interface{}]interface{}))
	}
}

func (c *ConfigSourceProvider) Root() MapNode {
	return MapNode{
		data:   c.data,
		parent: "",
	}
}

var (
	emptyMapNode = MapNode{}
)

func parentKey(parent, key string) string {
	if len(parent) > 0 {
		key = "." + key
	}
	return parent + key
}

type MapNode struct {
	change bool
	parent string
	data   map[interface{}]interface{}
}

func (m MapNode) Size() int {
	return len(m.data)
}

func (m MapNode) IsEmpty() bool {
	return 0 == m.Size()
}

func newMapNode(data map[interface{}]interface{}, parent, key string) MapNode {
	return MapNode{
		data:   data,
		parent: parentKey(parent, key),
	}
}

func (m MapNode) Keys() []string {
	keys := []string{}
	if len(m.data) > 0 {
		for k := range m.data {
			if v, ok := k.(string); ok {
				keys = append(keys, v)
			}
		}
	}
	return keys
}

func (m MapNode) Map(key string) MapNode {
	value, e := m.data[key]
	if e {
		return newMapNode(value.(map[interface{}]interface{}), m.parent, key)
	}
	if m.change {
		v := map[interface{}]interface{}{}
		m.data[key] = v
		return newMapNode(v, m.parent, key)
	}
	return emptyMapNode
}

func (m MapNode) String(key string, dv string) string {
	value, e := m.data[key]
	if e && value != nil {
		return value.(string)
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Int(key string, dv int) int {
	value, e := m.data[key]
	if e && value != nil {
		return value.(int)
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Float(key string, dv float64) float64 {
	value, e := m.data[key]
	if e && value != nil {
		return value.(float64)
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Bool(key string, dv bool) bool {
	value, e := m.data[key]
	if e && value != nil {
		return value.(bool)
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Duration(key string, dv time.Duration) time.Duration {
	value, e := m.data[key]
	if e && value != nil {
		tmp := value.(string)
		t, e := time.ParseDuration(tmp)
		if e == nil {
			return t
		}
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Time(key string, dv time.Time) time.Time {
	value, e := m.data[key]
	if e && value != nil {
		tmp := value.(string)
		t, e := time.Parse(time.RFC3339Nano, tmp)
		if e == nil {
			return t
		}
		t, e = time.Parse(time.RFC3339, tmp)
		if e == nil {
			return t
		}
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) StringL(key string, dv []string) []string {
	value, e := m.data[key]
	if e && value != nil {
		list := value.([]interface{})
		tmp := make([]string, len(list))
		for i, item := range list {
			tmp[i] = item.(string)
		}
		return tmp
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) IntL(key string, dv []int) []int {
	value, e := m.data[key]
	if e && value != nil {
		list := value.([]interface{})
		tmp := make([]int, len(list))
		for i, item := range list {
			tmp[i] = item.(int)
		}
		return tmp
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) FloatL(key string, dv []float64) []float64 {
	value, e := m.data[key]
	if e && value != nil {
		list := value.([]interface{})
		tmp := make([]float64, len(list))
		for i, item := range list {
			tmp[i] = item.(float64)
		}
		return tmp
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) BoolL(key string, dv []bool) []bool {
	value, e := m.data[key]
	if e && value != nil {
		list := value.([]interface{})
		tmp := make([]bool, len(list))
		for i, item := range list {
			tmp[i] = item.(bool)
		}
		return tmp
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}
