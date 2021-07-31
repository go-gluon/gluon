package utils

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]interface{}{})
}

func MapDeepCopy(m map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var mapCopy map[string]interface{}
	err = dec.Decode(&mapCopy)
	if err != nil {
		return nil, err
	}
	return mapCopy, nil
}
