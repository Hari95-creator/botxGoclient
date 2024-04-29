package dbconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	BusinessId    string `json:"Business-ID"`
	PhoneNumberId string `json:"Phone-Number-ID"`
	AccessToken   string `json:"User-Access-Token"`
	WabaId        string `json:"WABA-ID"`
	Version       string `json:"Version"`
	Url           string `json:"Url"`
}

// LoadConfig reads the configuration from config.json and returns a Config instance.
func LoadConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("Config file not found at %s", path)
		return nil, err
	}

	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Printf("Error unmarshalling config data: %v", err)
		return nil, err
	}

	return &config, nil
}
