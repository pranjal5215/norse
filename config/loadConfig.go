package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type jsonConfig map[string]interface{}


var (
	jsonConfigInstance jsonConfig
	loadErr error
	configFilePath = ""
)

func Configure(path string) {
	configFilePath = path
	jsonConfigInstance, loadErr = loadConfig()
	if loadErr != nil{
		os.Exit(1)
	}
}

// Load files from general config
func loadConfig() (jsonConfig, error) {
	if configFilePath == "" {
		fmt.Println("No config specified; pleae call config.Configure(filePath)")
		os.Exit(1)
	}
	file, err := ioutil.ReadFile(configFilePath)
	if err = json.Unmarshal(file, &jsonConfigInstance); err != nil {
		return nil, err
	}
	return jsonConfigInstance, err
}
