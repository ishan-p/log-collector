package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type ClientConfig struct {
	Watchers    []WatcherConfig `json:"watchers"`
	Collector   CollectorConfig `json:"server"`
	RetryParams RetryConfig     `json:"retry"`
}

func initiateClientConfigWithDefaults() ClientConfig {
	var config ClientConfig
	config.Collector = defaultServerConnection()
	config.RetryParams = defaultRetryConfiguration()
	return config
}

func validateClientConfig(config ClientConfig) (bool, error) {
	if len(config.Watchers) < 1 {
		return false, errors.New("Cannot initiate agent with 0 watchers")
	}
	// TODO: Add more validations around different values
	return true, nil
}

func ReadClientConfigJSON(configFilePath string) ClientConfig {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	isValidJson := json.Valid(data)
	if !isValidJson {
		log.Fatal("Provided config file is an invalid JSON")
	}

	config := initiateClientConfigWithDefaults()
	_ = json.Unmarshal(data, &config)
	isValid, validationError := validateClientConfig(config)
	if !isValid {
		log.Fatal(validationError)
	}
	return config
}

type ServerConfig struct {
	CollectorConfig
	Storage StorageConfig `json:"storage"`
}

func initiateServerConfigWithDefaults() ServerConfig {
	var config ServerConfig
	config.CollectorConfig = defaultServerConnection()
	return config
}

func ReadServerConfigJSON(configFilePath string) ServerConfig {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	isValidJson := json.Valid(data)
	if !isValidJson {
		log.Fatal("Provided config file is an invalid JSON")
	}

	config := initiateServerConfigWithDefaults()
	_ = json.Unmarshal(data, &config)
	return config
}
