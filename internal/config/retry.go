package config

type RetryConfig struct {
	MaxRetries   int `json:"max_retries"`
	MaxQueueSize int `json:"max_queue_size"`
}

func defaultRetryConfiguration() RetryConfig {
	var retryConfig RetryConfig
	retryConfig.MaxQueueSize = 500
	retryConfig.MaxRetries = 3
	return retryConfig
}
