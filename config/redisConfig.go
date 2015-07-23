package config

// Load redis configs after general config unmarshall
func loadRedisConfig() {
	temp = loadConfig()
	redisCs := make(map[string]map[string]string)
	for database, settings := range(temp["redis"].(map[string]interface{})){
		settingsI := settings.(map[string]interface{})
		configMap := make(map[string]string)
		for databaseConf, configMapI := range(settingsI){
			configMap[databaseConf] = configMapI.(string)
		}
		redisCs[database] = configMap
	}
	fmt.Println(redisCs["flight"]["port"])
}
