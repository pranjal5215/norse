package norse

import (
	"fmt"
)

dbname = "mysql"

// Load mysql configs after general config unmarshall
func LoadSqlConfig() (map[string]map[string]string, error) {
        // Get configVar from general config
        configVar, err := loadConfig()
        if err != nil{
                return nil, err
        }
        mysqlCs := make(map[string]map[string]string)
        for database, settings := range(configVar[dbname].(map[string]interface{})){
                settingsI := settings.(map[string]interface{})
                configMap := make(map[string]string)
                for databaseConf, configMapI := range(settingsI){
                        configMap[databaseConf] = configMapI.(string)
                }
                mysqlCs[database] = configMap
        }
        fmt.Println(mysqlCs["flight"]["port"])
        return mysqlCs, nil
}
/*
func main(){
	fmt.Println(LoadSqlConfig())
}*/
