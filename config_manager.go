package logcollector

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	Watchers    []Watcher          `json:"watchers"`
	Collector   ServerConnection   `json:"server"`
	RetryParams RetryConfiguration `json:"retry"`
}

func initiateConfigWithDefaults() Config {
	var config Config
	config.Collector = defaultServerConnection()
	config.RetryParams = defaultRetryConfiguration()
	return config
}

func validateConfig(config Config) (bool, error) {
	if len(config.Watchers) < 1 {
		return false, errors.New("Cannot initiate agent with 0 watchers")
	}
	// TODO: Add more validations around different values
	return true, nil
}

func readClientConfigJSON(configFilePath string) Config {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	isValidJson := json.Valid(data)
	if !isValidJson {
		log.Fatal("Provided config file is an invalid JSON")
	}

	config := initiateConfigWithDefaults()
	_ = json.Unmarshal(data, &config)
	isValid, validationError := validateConfig(config)
	if !isValid {
		log.Fatal(validationError)
	}
	return config
}
