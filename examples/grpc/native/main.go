package native

import (
	"fmt"
	"io/ioutil"
)

type KVType struct{}

func (KVType) Put(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin-native", string(value)))
	return ioutil.WriteFile("kv_"+key, value, 0644)
}

func (KVType) Get(key string) ([]byte, error) {
	return ioutil.ReadFile("kv_" + key)
}

func (KVType) Bench() string {
	return "1"
}

var KV KVType
