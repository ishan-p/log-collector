package logcollector

type Watcher struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}
