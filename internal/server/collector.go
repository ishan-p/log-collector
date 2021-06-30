package server

import (
	"encoding/json"
	"log"
	"net"

	"github.com/ishan-p/log-collector/internal/schema"
)

func handleCollect(conn net.Conn, storageConfig schema.StorageConfig) error {
	decoder := json.NewDecoder(conn)
	var collectPayload schema.CollectCmdPayload
	if err := decoder.Decode(&collectPayload); err != nil {
		log.Println("Failed to read the request")
		return err
	}
	ack := collect(collectPayload, storageConfig)
	collectResponse := schema.CollectCmdResponse{Ack: ack}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&collectResponse); err != nil {
		log.Println("Failed to send the ack")
		return err
	}
	return nil
}
