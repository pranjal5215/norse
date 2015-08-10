package config

// Load redis configs after general config unmarshall
func LoadRedisConfig() (map[string]map[string]string, error) {
	var dbname = "redis"
	redisCs := make(map[string]map[string]string)
	// Convert package conf jsonConfigInstance to individual vertical config.
	for database, settings := range jsonConfigInstance[dbname].(map[string]interface{}) {
		settingsI := settings.(map[string]interface{})
		configMap := make(map[string]string)
		for databaseConf, configMapI := range settingsI {
			configMap[databaseConf] = configMapI.(string)
		}
		redisCs[database] = configMap
	}
	return redisCs, nil
}
