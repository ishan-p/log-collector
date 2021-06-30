package config

type WatcherConfig struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}
