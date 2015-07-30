package norse

import (
	"fmt"
)

// Load redis configs after general config unmarshall
func LoadRedisConfig() (map[string]map[string]string, error) {
	var dbname = "redis"
	// Get configVar from general config
	configVar, err := loadConfig()
	if err != nil {
		return nil, err
	}
	redisCs := make(map[string]map[string]string)
	for database, settings := range configVar[dbname].(map[string]interface{}) {
		settingsI := settings.(map[string]interface{})
		configMap := make(map[string]string)
		for databaseConf, configMapI := range settingsI {
			configMap[databaseConf] = configMapI.(string)
		}
		redisCs[database] = configMap
	}
	fmt.Println(redisCs["flight"]["port"])
	return redisCs, nil
}
