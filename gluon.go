package gluon

import (
	"embed"
	"fmt"
	"time"

	"github.com/go-gluon/gluon/config"
	"github.com/go-gluon/gluon/log"
)

var (
	G     *Gluon
	start = time.Now()
)

type Annotation struct{}

type GluonInfo struct {
	appName    string
	appVersion string
	version    string
}

func (a GluonInfo) AppName() string {
	return a.appName
}

func (a GluonInfo) AppVersion() string {
	return a.appVersion
}

func (a GluonInfo) Version() string {
	return a.version
}

type Runtime struct {
	config    *config.ConfigSourceProvider
	resources *embed.FS
}

func (r *Runtime) init() {
	// core modules
	if r.resources != nil {
		err := r.config.Init(r.resources)
		if err != nil {
			panic(err)
		}
	}
}

type Gluon struct {
	info     *GluonInfo
	runtime  *Runtime
	ext      []*ExtensionRuntime
	services []*ExtensionRuntime
}

func CreateGluon(appName, appVersion, version string, resources *embed.FS, ext ...*ExtensionRuntime) *Gluon {

	G = &Gluon{
		info:     &GluonInfo{appName: appName, appVersion: appVersion, version: version},
		runtime:  &Runtime{resources: resources, config: config.Default},
		ext:      ext,
		services: []*ExtensionRuntime{},
	}
	if err := G.init(); err != nil {
		panic(err)
	}
	return G
}

func (g *Gluon) init() error {

	// initialize runtime core modules
	g.runtime.init()

	// activate extension
	size := len(g.ext)
	names := make([]string, size)
	if size > 0 {

		// register extension
		for i, e := range g.ext {
			err := g.initExtension(e)
			if err != nil {
				panic(err)
			}
			names[i] = e.name

			if e.service {
				g.services = append(g.services, e)
			}
		}
	}

	log.Info(fmt.Sprintf("gluon extensions %v", names))
	return nil
}

func (r *Gluon) initExtension(e *ExtensionRuntime) error {

	// setup extension configuration
	econfig := e.ext.InitConfig()
	if econfig != nil {
		if reader, ok := econfig.(config.ConfigReader); ok {
			if err := config.Default.InitExtension(e.name, reader); err != nil {
				return err
			}
		} else {
			log.Warn("Configuration does not implements ConfigReader", log.Fields{"extension": e.name})
		}
	}

	// initialize the extension
	return e.ext.Init(r.info, r.runtime)
}

func (g *Gluon) Start() {
	shutdown := make(chan bool)
	for _, e := range g.services {
		x := e
		go func() {
			x.ext.Start()
			shutdown <- true
		}()
	}
	log.Info(fmt.Sprintf("%v %v (powered by gluon %v) started in %v", g.info.appName, g.info.appVersion, g.info.version, time.Since(start)))
	<-shutdown
}
