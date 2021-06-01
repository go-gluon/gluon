package config

import (
	"os"
	"regexp"
	"strings"
)

var envRegexp = createRegexp()

// initialize yaml file configuration source
func createRegexp() *regexp.Regexp {

	// initialize the regex for EnvConfigSource
	tmp, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	return tmp
}

type EnvConfigSource struct {
	envs map[string]string
}

func (f *EnvConfigSource) Init() error {
	f.envs = map[string]string{}
	return f.parseEnv()
}

func (f *EnvConfigSource) Priority() int {
	return 300
}

func (f *EnvConfigSource) Name() string {
	return `env`
}

func (f *EnvConfigSource) Property(name string) (string, bool, error) {
	tmp := envRegexp.ReplaceAllString(name, "_")
	tmp = strings.ToUpper(tmp)
	v, e := f.envs[tmp]
	return v, e, nil
}

func (f *EnvConfigSource) Properties() (map[string]string, error) {
	return f.envs, nil
}

func (f *EnvConfigSource) parseEnv() error {
	items := os.Environ()
	if len(items) > 0 {
		for _, item := range items {
			tmp := strings.Split(item, "=")
			f.envs[tmp[0]] = tmp[1]
		}
	}
	return nil
}
