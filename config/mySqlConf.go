package config

import (
	"fmt"
)

// Load mysql configs after general config unmarshall
func LoadSqlConfig() (map[string]map[string]string, error) {
	var dbname = "mysql"
	mysqlCs := make(map[string]map[string]string)
	for database, settings := range jsonConfigInstance[dbname].(map[string]interface{}) {
		settingsI := settings.(map[string]interface{})
		configMap := make(map[string]string)
		for databaseConf, configMapI := range settingsI {
			configMap[databaseConf] = configMapI.(string)
		}
		mysqlCs[database] = configMap
	}
	fmt.Println(mysqlCs["flight"]["port"])
	return mysqlCs, nil
}
