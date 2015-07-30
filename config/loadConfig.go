package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type jsonConfig map[string]interface{}

var configFilePath = ""

func Configure(path string) {
	configFilePath = path
}

// Load files from general config
func loadConfig() (jsonConfig, error) {
	if configFilePath == "" {
		fmt.Println("No config specified; pleae call config.Configure(filePath)")
		os.Exit(1)
	}
	file, err := ioutil.ReadFile(configFilePath)
	temp := make(jsonConfig)
	if err = json.Unmarshal(file, &temp); err != nil {
		return nil, err
	}
	return temp, err
}
