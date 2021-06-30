package schema

type ClientConfig struct {
	Watchers    []WatcherConfig `json:"watchers"`
	Collector   CollectorConfig `json:"server"`
	RetryParams RetryConfig     `json:"retry"`
}

type WatcherConfig struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}

type RetryConfig struct {
	MaxRetries   int `json:"max_retries"`
	MaxQueueSize int `json:"max_queue_size"`
}

type CollectorConfig struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	ServerWaitTimeSec int    `json:"wait_time"`
}
