package deployer

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ServerConfig struct {
	Port int `json:"port"`
}

func LoadConfig(path string, config interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error loading "+path+": ", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("Error parsing JSON from "+path+": ", err)
	}
}
