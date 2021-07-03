package server

import (
	"encoding/json"
	"log"
	"net"

	"github.com/ishan-p/log-collector/internal/schema"
)

type CollectCommand schema.CollectRequest

func (cmd CollectCommand) Execute(handler RequestHandler) (ExecutionStatus, error) {
	storer, err := NewStorer(cmd.Destination, handler.Config.Storage)
	if err != nil {
		log.Println(err)
		return false, err
	}
	jsonLogEvent, err := json.Marshal(cmd)
	if err != nil {
		log.Println("Unable to encode log as json")
		return false, err
	}
	status, err := storer.write(jsonLogEvent)
	if err != nil {
		log.Println(err)
		return ExecutionStatus(status), err
	}
	return ExecutionStatus(status), nil
}

func (cmd CollectCommand) Reply(conn net.Conn, status ExecutionStatus) error {
	response := schema.CollectResponse{Ack: bool(status)}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&response); err != nil {
		log.Println("Failed to send the ack")
		return err
	}
	return nil
}
