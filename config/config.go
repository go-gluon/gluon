package config

import (
	"io/fs"
	"strconv"
	"strings"
	"time"
)

const (
	configPrefix = "gluon"
)

var (
	// Default configuration source provider instance
	Default *ConfigSourceProvider
	// ConfigProfileProperty configuration profile property
	configProfileProperty = configPrefix + ".config.profile"
)

type ConfigReader interface {
	ReadFromMapNode(node MapNode) error
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
		cache:   map[string]void{},
	}

	// set profile
	if profile, e := Default.getSourceRawValue(configProfileProperty); e {
		Default.SetProfile(profile)
	}
}

// SetProfile set profile to default provider
func SetProfile(profile string) {
	Default.SetProfile(profile)
}

// Profile of the default configuration source provider
func Profile() string {
	return Default.profile
}

type void struct{}

// ConfigSourceProvider configuration source provider
type ConfigSourceProvider struct {
	sources      []ConfigSource
	data         map[interface{}]interface{}
	profile      string
	cache        map[string]void
	mapNode      MapNode
	extensionMap MapNode
}

func (c *ConfigSourceProvider) InitExtension(name string, reader ConfigReader) error {
	return reader.ReadFromMapNode(c.extensionMap.Map(name))
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

func (c *ConfigSourceProvider) isProfile() bool {
	return len(c.profile) > 0
}

func (c *ConfigSourceProvider) getSourceValue(parent, key string) (string, bool) {
	key = parent + "." + key
	if _, e := c.cache[key]; e {
		return "", false
	}

	c.cache[key] = void{}

	for _, s := range c.sources {
		if c.isProfile() {
			if v, e, _ := s.GetRawValue(c.profile + "." + key); e {
				return v, true
			}
		}
		if v, e, _ := s.GetRawValue(key); e {
			return v, true
		}
	}
	return "", false
}

func (c *ConfigSourceProvider) getSourceRawValue(name string) (string, bool) {
	for _, s := range c.sources {
		if v, e, _ := s.GetRawValue(name); e {
			return v, true
		}
	}
	return "", false
}

func (c *ConfigSourceProvider) Init(resources fs.FS) error {
	data, err := loadYaml(resources)
	if err != nil {
		return err
	}
	c.data = data

	// merge profile
	if c.isProfile() {
		if d, e := c.data[c.profile]; e {
			if data, ee := d.(map[interface{}]interface{}); ee {
				merge(data, c.data)
			}
		}
	}

	c.mapNode = newMapNode(c, c.data, "", "")
	c.extensionMap = c.mapNode.Map(configPrefix)
	c.extensionMap.change = true

	return nil
}

func merge(profile, config map[interface{}]interface{}) {

	for key, value := range profile {

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

func (c *ConfigSourceProvider) Map() MapNode {
	return c.mapNode
}

func (c *ConfigSourceProvider) String(key string, dv string) string {
	m, k := c.findMap(key)
	return m.String(k, dv)
}

func (c *ConfigSourceProvider) Int(key string, dv int) int {
	m, k := c.findMap(key)
	return m.Int(k, dv)
}

func (c *ConfigSourceProvider) Float(key string, dv float64) float64 {
	m, k := c.findMap(key)
	return m.Float64(k, dv)
}

func (c *ConfigSourceProvider) Bool(key string, dv bool) bool {
	m, k := c.findMap(key)
	return m.Bool(k, dv)
}

func (c *ConfigSourceProvider) Time(key string, dv time.Time) time.Time {
	m, k := c.findMap(key)
	return m.Time(k, dv)
}

func (c *ConfigSourceProvider) Duration(key string, dv time.Duration) time.Duration {
	m, k := c.findMap(key)
	return m.Duration(k, dv)
}

func (c *ConfigSourceProvider) findMap(key string) (MapNode, string) {
	if len(key) == 0 {
		return MapNode{}, key
	}
	items := strings.Split(key, ".")

	n := c.Map()
	var i int
	for i = 0; i < len(items)-1 && !n.IsEmpty(); i++ {
		n = n.Map(items[i])
	}
	return n, items[len(items)-1]
}

var emptyMapNode = MapNode{}

func parentKey(parent, key string) string {
	if len(parent) > 0 {
		key = "." + key
	}
	return parent + key
}

type MapNode struct {
	provider *ConfigSourceProvider
	change   bool
	parent   string
	data     map[interface{}]interface{}
}

func (m MapNode) Size() int {
	return len(m.data)
}

func (m MapNode) IsEmpty() bool {
	return 0 == m.Size()
}

func newMapNode(c *ConfigSourceProvider, data map[interface{}]interface{}, parent, key string) MapNode {
	return MapNode{
		provider: c,
		data:     data,
		parent:   parentKey(parent, key),
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
	if value, e := m.data[key]; e {
		return newMapNode(m.provider, value.(map[interface{}]interface{}), m.parent, key)
	}
	if m.change {
		v := map[interface{}]interface{}{}
		m.data[key] = v
		return newMapNode(m.provider, v, m.parent, key)
	}
	return emptyMapNode
}

func (m MapNode) String(key string, dv string) string {

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		m.data[key] = sv
		return sv
	}

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

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		if v, e := strconv.Atoi(sv); e == nil {
			m.data[key] = v
			return v
		}
	}

	value, e := m.data[key]
	if e && value != nil {
		return value.(int)
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Float64(key string, dv float64) float64 {

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		if v, e := strconv.ParseFloat(sv, 64); e == nil {
			m.data[key] = v
			return v
		}
	}

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

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		if v, e := strconv.ParseBool(sv); e == nil {
			m.data[key] = v
			return v
		}
	}

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

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		if v, e := time.ParseDuration(sv); e == nil {
			m.data[key] = v
			return v
		}
	}

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

	if sv, se := m.provider.getSourceValue(m.parent, key); se {
		if v, e := stringToTime(sv); e == nil {
			m.data[key] = v
			return v
		}
	}

	value, e := m.data[key]
	if e && value != nil {
		tmp := value.(string)
		t, e := stringToTime(tmp)
		if e == nil {
			return t
		}
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func stringToTime(tmp string) (time.Time, error) {
	if t, e := time.Parse(time.RFC3339Nano, tmp); e == nil {
		return t, nil
	}
	t, e := time.Parse(time.RFC3339, tmp)
	if e == nil {
		return t, nil
	}
	return time.Time{}, e
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

func (m MapNode) Float64L(key string, dv []float64) []float64 {
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

func (m MapNode) StringM(key string, dv map[string]string) map[string]string {
	value, e := m.data[key]
	if !e || value == nil {
		return dv
	}
	if m, ok := value.(map[interface{}]interface{}); ok {
		r := map[string]string{}
		for k, v := range m {
			r[k.(string)] = v.(string)
		}
		return r
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) IntM(key string, dv map[string]int) map[string]int {
	value, e := m.data[key]
	if !e || value == nil {
		return dv
	}
	if m, ok := value.(map[interface{}]interface{}); ok {
		r := map[string]int{}
		for k, v := range m {
			r[k.(string)] = v.(int)
		}
		return r
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) Float64M(key string, dv map[string]float64) map[string]float64 {
	value, e := m.data[key]
	if !e || value == nil {
		return dv
	}
	if m, ok := value.(map[interface{}]interface{}); ok {
		r := map[string]float64{}
		for k, v := range m {
			r[k.(string)] = v.(float64)
		}
		return r
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}

func (m MapNode) BoolM(key string, dv map[string]bool) map[string]bool {
	value, e := m.data[key]
	if !e || value == nil {
		return dv
	}
	if m, ok := value.(map[interface{}]interface{}); ok {
		r := map[string]bool{}
		for k, v := range m {
			r[k.(string)] = v.(bool)
		}
		return r
	}
	if m.change {
		m.data[key] = dv
	}
	return dv
}
