package schema

type CommandConatiner map[string]interface{}

var CollectCmd string

type CollectRequest struct {
	Timestamp   int64               `json:"timestamp"`
	Tags        []map[string]string `json:"tags"`
	Record      string              `json:"record"`
	Destination string              `json:"destination"`
}

type CollectResponse struct {
	Ack bool `json:"ack"`
}

func init() {
	CollectCmd = "collect"
}
