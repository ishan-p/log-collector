package logcollector

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
