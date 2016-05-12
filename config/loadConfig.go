package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var IsRedisConfigured bool
var IsMemcacheConfigured bool
var IsMysqlConfigured bool

type jsonConfig map[string]interface{}

var (
	jsonConfigInstance jsonConfig
	loadErr            error
	configFilePath     = ""
)

func Configure(path string) {
	configFilePath = path
	jsonConfigInstance, loadErr = loadConfig()
	if loadErr != nil {
		panic("ERROR LOADING CONFIG: " + loadErr.Error())
	}
	_, IsRedisConfigured = jsonConfigInstance["redis"]
	_, IsMemcacheConfigured = jsonConfigInstance["memcache"]
	_, IsMysqlConfigured = jsonConfigInstance["mysql"]
}

// Load files from general config
func loadConfig() (jsonConfig, error) {
	if configFilePath == "" {
		fmt.Println("No config specified; pleae call config.Configure(filePath)")
		os.Exit(1)
	}
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(file, &jsonConfigInstance); err != nil {
		return nil, err
	}
	return jsonConfigInstance, nil
}
