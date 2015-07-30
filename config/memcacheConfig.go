package config

import (
	"strings"
)

func LoadMemcacheConfig() (map[string][]string, error) {
	var dbname = "memcache"
	// Get configVar from general config
	configVar, err := loadConfig()
	if err != nil {
		return nil, err
	}
	memConfig := make(map[string][]string)
	for memType, memPortStr := range configVar[dbname].(map[string]interface{}) {
		strLst := strings.Split(memPortStr.(string), ",")
		memConfig[memType] = strLst
	}
	return memConfig, nil
}
