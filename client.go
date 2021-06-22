package logcollector

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

func RunAgent(configFile string) {
	config := readClientConfigJSON(configFile)
	done := make(chan bool)
	for _, watcher := range config.Watchers {
		log.Printf("Watching file %s\n", watcher.FileName)
		notificationChannel := make(chan string)
		defer close(notificationChannel)
		go watch(watcher.FileName, notificationChannel)
		go notify(notificationChannel, config.Collector, watcher.Destination, watcher.Tags, config.RetryParams)
	}
	<-done
}

func notify(ch chan string, server ServerConnection, destination string, tags []map[string]string, retryConfig RetryConfiguration) {
	for notification := range ch {
		sendToServer(notification, server, destination, tags, retryConfig)
	}
}

func sendToServer(notification string, server ServerConnection, destination string, tags []map[string]string, retryConfig RetryConfiguration) {
	dest := server.Host + ":" + strconv.Itoa(server.Port)
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Println("Could not connect to the server")
	}
	if conn == nil {
		log.Println("Failed to create successful connection")
	}
	defer conn.Close()

	commandResponse, err := initCollectRequest(conn)
	if err != nil || !commandResponse.Begin {
		log.Println("Failed to initiate command")
	}

	sendLog(conn, notification, destination, tags)
}

func initCollectRequest(conn net.Conn) (CommandResponse, error) {
	var commandResponse CommandResponse
	commandRequest := CommandRequest{
		Command: "collect",
	}
	if err := json.NewEncoder(conn).Encode(&commandRequest); err != nil {
		log.Println(err)
	}
	err := json.NewDecoder(conn).Decode(&commandResponse)
	return commandResponse, err
}

func sendLog(conn net.Conn, notification string, destination string, tags []map[string]string) {
	collectCmdPayload := CollectCmdPayload{
		Record:      notification,
		Timestamp:   time.Now().Unix(),
		Tags:        tags,
		Destination: destination,
	}
	if err := json.NewEncoder(conn).Encode(&collectCmdPayload); err != nil {
		log.Println(err)
	}
	var collectResp CollectCmdResponse
	json.NewDecoder(conn).Decode(&collectResp)
	log.Println(collectResp)
}
