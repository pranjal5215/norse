package config

func loadConfig(){
	file, err := ioutil.ReadFile(configFilePath)
	temp := make(JsonConfig)
	if err = json.Unmarshal(file, &temp); err != nil {
		fmt.Println("Error in parsing config: ", err)
	}
}

