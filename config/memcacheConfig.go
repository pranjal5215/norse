
package norse

import (
	"fmt"
)

var dbname = "memcache"

func LoadMemcacheConfig() (map[string][]string, error){
	// Get configVar from general config
	configVar, err := loadConfig()
	if err != nil{
		return nil, err
	}
	memConfig := make(map[string][]string)
	for memType, memPortStr := range(configVar[dbname].(map[string]interface{})){
		strLst := strings.Split(memPortStr, ",")
		memConfig[memType] = strLst
	}
	return memConfig, nil
}
