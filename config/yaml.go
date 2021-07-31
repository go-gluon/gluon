package config

import (
	"errors"
	"io/fs"

	"gopkg.in/yaml.v2"
)

const resourceFile = "application.yaml"

var errorFind = errors.New("find item")

func loadYaml(resources fs.FS) (map[interface{}]interface{}, error) {
	data := map[interface{}]interface{}{}

	rf := ""
	er := fs.WalkDir(resources, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if d.Name() == resourceFile {
				rf = path
				return errorFind
			}
		}
		return nil
	})
	if er != nil && er != errorFind {
		return data, er
	}
	if len(rf) == 0 {
		return data, nil
	}

	d, err := fs.ReadFile(resources, rf)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(d, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
