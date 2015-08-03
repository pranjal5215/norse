package config

import (
	"strings"
)

func LoadMemcacheConfig() (map[string][]string, error) {
	var dbname = "memcache"
	memConfig := make(map[string][]string)
	for memType, memPortStr := range jsonConfigInstance[dbname].(map[string]interface{}) {
		strLst := strings.Split(memPortStr.(string), ",")
		memConfig[memType] = strLst
	}
	return memConfig, nil
}
