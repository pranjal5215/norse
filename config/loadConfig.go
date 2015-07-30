package norse

import (
	"encoding/json"
	"io/ioutil"
)

type jsonConfig map[string]interface{}

var configFilePath = "/home/pranjal/goOld/config.json"

// Load files from general config
func loadConfig() (jsonConfig, error) {
	file, err := ioutil.ReadFile(configFilePath)
	temp := make(jsonConfig)
	if err = json.Unmarshal(file, &temp); err != nil {
		return nil, err
	}
	return temp, err
}
