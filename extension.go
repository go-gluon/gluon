package gluon

type ExtensionRuntime struct {
	name     string
	priority int
	version  string
	service  bool
	ext      Extension
}

func (e *ExtensionRuntime) Extension() Extension {
	return e.ext
}

func NewExtension(name string, priority int, version string, service bool, ext Extension) *ExtensionRuntime {
	return &ExtensionRuntime{
		name:     name,
		priority: priority,
		version:  version,
		service:  service,
		ext:      ext,
	}
}

type Extension interface {
	InitConfig() interface{}
	Init(info *GluonInfo, runtime *Runtime) error
	Start()
}
