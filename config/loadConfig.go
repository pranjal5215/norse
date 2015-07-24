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
		fmt.Println("Error in parsing config: ", err)
		return nil, err
	}
	return temp, err
}

