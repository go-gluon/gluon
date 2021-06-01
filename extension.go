package gluon

import (
	"embed"
	"sort"

	"github.com/go-gluon/gluon/config"
	"github.com/go-gluon/gluon/log"
)

type ExtensionInit = func(resources embed.FS, config interface{}) error

type Extension struct {
	Name     string
	Priority int
	Init     ExtensionInit
	Config   interface{}
}

type ExtensionProvider interface {
	NewExtesion() Extension
}

func RegisterExtensions(resources embed.FS, providers ...ExtensionProvider) error {
	// core modules
	err := config.AddYaml(resources)
	if err != nil {
		panic(err)
	}

	// activate extension
	if len(providers) > 0 {
		extensions := make([]Extension, len(providers))
		for i, e := range providers {
			extensions[i] = e.NewExtesion()
		}

		// sort extension base on the priority
		sort.Slice(extensions, func(i, j int) bool {
			return extensions[i].Priority < extensions[j].Priority
		})

		tmp := make([]string, len(extensions))
		for i, e := range extensions {
			tmp[i] = e.Name

			err := config.Extension(e.Name, e.Config)
			if err != nil {
				return err
			}

			err = e.Init(resources, e.Config)
			if err != nil {
				panic(err)
			}
		}
		log.Info("Loaded extension", log.Fields{"extensions": tmp})
	}
	return nil
}
