package gluon

import (
	"runtime/debug"

	"github.com/go-gluon/gluon/log"
)

func version() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "0.0.0"
	}
	for _, m := range bi.Deps {
		log.Debug("Deps", log.Fields{"path": m.Path, "version": m.Version})
	}
	return bi.Main.Version
}
