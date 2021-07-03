package server

import (
	"errors"
	"net"
)

type ExecutionStatus bool

type Command interface {
	Execute(handler RequestHandler) (ExecutionStatus, error)
	Reply(conn net.Conn, status ExecutionStatus) error
}

var InvalidCmdFormatErr, MissingCmdKeyErr, InvalidCmdErr, MissingPayloadKeyErr, InvalidPayloadErr error

func init() {
	InvalidCmdFormatErr = errors.New("Invalid command format")
	MissingCmdKeyErr = errors.New("Missing key 'cmd' in the request")
	InvalidCmdErr = errors.New("Invalid command")
	MissingPayloadKeyErr = errors.New("Missing key 'payload' in the request")
	InvalidPayloadErr = errors.New("Invalid request payload")
}

func NewCommand(cmd interface{}) (Command, error) {
	cmdVal, ok := cmd.(map[string]interface{})
	if !ok {
		return nil, InvalidCmdFormatErr
	}
	if cmdVal["cmd"] == nil {
		return nil, MissingCmdKeyErr
	}
	if cmdVal["payload"] == nil {
		return nil, MissingPayloadKeyErr
	}

	switch c := cmdVal["cmd"]; c {
	case "collect":
		var command CollectCommand
		commandData, ok := cmdVal["payload"].(map[string]interface{})
		if !ok {
			return nil, InvalidPayloadErr
		}
		command.Destination = commandData["destination"].(string)
		command.Record = commandData["record"].(string)
		for _, t := range commandData["tags"].([]interface{}) {
			tagI := t.(map[string]interface{})
			tag := make(map[string]string)
			for k, v := range tagI {
				tag[k] = v.(string)
			}
			command.Tags = append(command.Tags, tag)
		}
		command.Timestamp = int64(commandData["timestamp"].(float64))
		return command, nil
	default:
		return nil, InvalidCmdErr
	}
}
