package server

import (
	"encoding/json"
	"log"
	"net"

	serverConfig "github.com/ishan-p/log-collector/internal/config"
)

func handleCollect(conn net.Conn, storageConfig serverConfig.StorageConfig) error {
	decoder := json.NewDecoder(conn)
	var collectPayload serverConfig.CollectCmdPayload
	if err := decoder.Decode(&collectPayload); err != nil {
		log.Println("Failed to read the request")
		return err
	}
	ack := collect(collectPayload, storageConfig)
	collectResponse := serverConfig.CollectCmdResponse{Ack: ack}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&collectResponse); err != nil {
		log.Println("Failed to send the ack")
		return err
	}
	return nil
}
