package server

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ishan-p/log-collector/internal/schema"
)

func initiateServerConfigWithDefaults() schema.ServerConfig {
	var config schema.ServerConfig
	config = defaultServerConnection()
	return config
}

func ReadJSON(configFilePath string) schema.ServerConfig {
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

func defaultServerConnection() schema.ServerConfig {
	var connection schema.ServerConfig
	connection.Host = "127.0.0.1"
	connection.Port = 8888
	connection.ServerWaitTimeSec = 3
	connection.ConnectionIdleTime = 45 * time.Second
	connection.MaxConnectionRetries = 5
	connection.SleepRetryDuration = 10 * time.Millisecond
	return connection
}
