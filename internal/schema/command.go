package schema

type CollectCmdPayload struct {
	Timestamp   int64               `json:"timestamp"`
	Tags        []map[string]string `json:"tags"`
	Record      string              `json:"record"`
	Destination string              `json:"destination"`
}

type CollectCmdResponse struct {
	Ack bool `json:"ack"`
}

type CommandRequest struct {
	Command string `json:"cmd"`
}

type CommandResponse struct {
	Command string `json:"cmd"`
	Begin   bool   `json:"begin"`
}
