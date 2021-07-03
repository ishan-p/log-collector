package client

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/ishan-p/log-collector/internal/schema"
)

func initiateClientConfigWithDefaults() schema.ClientConfig {
	var config schema.ClientConfig
	config.Collector = defaultServerConnection()
	config.RetryParams = defaultRetryConfiguration()
	return config
}

func validateClientConfig(config schema.ClientConfig) (bool, error) {
	if len(config.Watchers) < 1 {
		return false, errors.New("Cannot initiate agent with 0 watchers")
	}
	// TODO: Add more validations around different values
	return true, nil
}

func ReadJSON(configFilePath string) schema.ClientConfig {
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

func defaultRetryConfiguration() schema.RetryConfig {
	var retryConfig schema.RetryConfig
	retryConfig.MaxQueueSize = 500
	retryConfig.MaxRetries = 3
	return retryConfig
}

func defaultServerConnection() schema.CollectorConfig {
	var connection schema.CollectorConfig
	connection.Host = "127.0.0.1"
	connection.Port = 8888
	connection.ServerWaitTimeSec = 3
	return connection
}
