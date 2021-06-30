package config

import "time"

type CollectorConfig struct {
	Host                 string        `json:"host"`
	Port                 int           `json:"port"`
	ServerWaitTimeSec    int           `json:"wait_time"`
	ConnectionIdleTime   time.Duration `json:"connection_idle_time"`
	MaxConnectionRetries int           `json:"max_connection_retries"`
	SleepRetryDuration   time.Duration `json:"sleep_retry"`
}

func defaultServerConnection() CollectorConfig {
	var connection CollectorConfig
	connection.Host = "127.0.0.1"
	connection.Port = 8888
	connection.ServerWaitTimeSec = 3
	connection.ConnectionIdleTime = 45 * time.Second
	connection.MaxConnectionRetries = 5
	connection.SleepRetryDuration = 10 * time.Millisecond
	return connection
}
