package gluon

import (
	"embed"
	"sort"

	"github.com/go-gluon/gluon/config"
	"github.com/go-gluon/gluon/log"
)

type ExtensionInfo struct {
	Name     string
	Priority int
}

type Extension interface {
	Info() ExtensionInfo
	Configuration() interface{}
	Create(resources *embed.FS) error
}

type ExtensionProvider interface {
	NewExtension() Extension
}

var extensions []Extension

func RegisterExtensions(resources *embed.FS, providers ...ExtensionProvider) error {
	// core modules
	if resources != nil {
		err := config.Default.Init(resources)
		if err != nil {
			panic(err)
		}
	}

	// activate extension
	names := []string{}
	if len(providers) > 0 {
		extensions = make([]Extension, len(providers))
		tmp := make([]string, len(providers))

		for i, e := range providers {
			ex := e.NewExtension()
			extensions[i] = ex
			tmp[i] = ex.Info().Name
		}

		// sort extension base on the priority
		sort.Slice(extensions, func(i, j int) bool {
			return extensions[i].Info().Priority < extensions[j].Info().Priority
		})

		// register extension
		for _, e := range extensions {
			err := registerExtensions(e, resources)
			if err != nil {
				panic(err)
			}
		}

		names = append(names, tmp...)
	}

	log.Info("Gluon", log.Fields{"version": "v0.0.0", "extensions": names})
	return nil
}

func registerExtensions(e Extension, resources *embed.FS) error {
	// if e.Configuration() != nil {
	// 	// err := config.Extension(e.Info().Name, e.Configuration())
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }
	return e.Create(resources)
}
