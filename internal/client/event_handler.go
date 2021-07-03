package client

import (
	"log"

	"github.com/ishan-p/log-collector/internal/platform/tcp"
	"github.com/ishan-p/log-collector/internal/schema"
)

type ServerConnection struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	ServerWaitTimeSec int    `json:"wait_time"`
}

func sendToServer(notification Notification, server schema.CollectorConfig, retryChannel chan Notification) {
	client := tcp.NewClient(server.Host, server.Port)
	conn, err := client.Connect()
	if err != nil {
		log.Println("Failed to create successful connection")
		retryChannel <- notification
		return
	} else {
		defer conn.Close()
	}
	command := buildCollectCommand(notification)
	response, err := client.Send(conn, command)
	if err != nil {
		retryChannel <- notification
		return
	}
	ackI, ok := response.(map[string]interface{})
	if !ok {
		log.Println("Could not decode response")
		retryChannel <- notification
		return
	}
	ack := schema.CollectResponse{}
	ack.Ack = ackI["ack"].(bool)
	if !ack.Ack {
		log.Println("Did not receive log acknowledgement")
		retryChannel <- notification
		return
	}
}

func buildCollectCommand(notification Notification) schema.CommandConatiner {
	command := make(schema.CommandConatiner)
	command["cmd"] = schema.CollectCmd
	collectCmdPayload := schema.CollectRequest{
		Record:      notification.LogEvent,
		Timestamp:   notification.Timestamp,
		Tags:        notification.Tags,
		Destination: notification.Destination,
	}
	command["payload"] = collectCmdPayload
	return command
}
