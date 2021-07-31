package config

import (
	"errors"
	"os"
	"strings"
)

type FlagsConfigSource struct {
	args  []string
	flags map[string]string
}

func (f *FlagsConfigSource) Init() error {
	f.flags = map[string]string{}
	f.args = os.Args[1:]
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}
		return err
	}
	return nil
}

func (f *FlagsConfigSource) GetRawValue(name string) (string, bool, error) {
	tmp := strings.ReplaceAll(name, ".", "-")
	v, e := f.flags[tmp]
	return v, e, nil
}

func (f *FlagsConfigSource) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}
	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, errors.New("bad flag syntax: " + s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	// It must have a value, which might be the next argument.
	if !hasValue && len(f.args) > 0 {
		// value is the next arg
		tmp := f.args[0]
		if len(tmp) > 0 && tmp[0] != '-' {
			value = tmp
			f.args = f.args[1:]
		}
	}
	f.flags[name] = value

	return true, nil
}
