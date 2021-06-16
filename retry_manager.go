package logcollector

type RetryConfiguration struct {
	MaxRetries   int `json:"max_retries"`
	MaxQueueSize int `json:"max_queue_size"`
}

func defaultRetryConfiguration() RetryConfiguration {
	var retryConfig RetryConfiguration
	retryConfig.MaxQueueSize = 500
	retryConfig.MaxRetries = 3
	return retryConfig
}
