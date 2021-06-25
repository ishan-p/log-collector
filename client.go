package logcollector

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

type Notification struct {
	LogEvent     string
	Timestamp    int64
	Tags         []map[string]string
	Destination  string
	RetryAttempt int
	LastRetry    int64
}

func RunAgent(configFile string) {
	config := readClientConfigJSON(configFile)
	done := make(chan bool)
	for _, watcher := range config.Watchers {
		log.Printf("Watching file %s\n", watcher.FileName)
		watcher.Tags = append(watcher.Tags, map[string]string{"source_file": watcher.FileName})
		notificationChannel := make(chan string)
		defer close(notificationChannel)
		go watch(watcher.FileName, notificationChannel)
		go notify(notificationChannel, config.Collector, watcher.Destination, watcher.Tags)
	}
	<-done
}

func notify(ch chan string, server ServerConnection, destination string, tags []map[string]string) {
	for event := range ch {
		notification := Notification{
			LogEvent:     event,
			Timestamp:    time.Now().Unix(),
			Tags:         tags,
			Destination:  destination,
			RetryAttempt: 0,
			LastRetry:    time.Now().Unix(),
		}
		sendToServer(notification, server)
	}
}

func sendToServer(notification Notification, server ServerConnection) {
	conn := initServerConnection(server.Host, server.Port)
	if conn == nil {
		log.Println("Failed to create successful connection")
	}
	defer conn.Close()

	commandResponse, err := initCollectRequest(conn)
	if err != nil || !commandResponse.Begin {
		log.Fatalln("Failed to initiate command: ", err)
	}

	collectAck, err := sendLog(conn, notification)
	if err != nil {
		log.Fatalln("Failed to send log event: ", err)
	}
	if !collectAck.Ack {
		log.Println("Did not receive log acknowledgement")
	}
}

func initServerConnection(host string, port int) net.Conn {
	var connectionRetryCount int
	var conn net.Conn
	var err error
	dialer := &net.Dialer{
		Timeout:   time.Second * 300,
		KeepAlive: time.Minute * 5,
	}
	maxConnectionRetries := 5
	defaultRetrySleep := time.Second * 1
	connectionRetrySleep := defaultRetrySleep
	dest := host + ":" + strconv.Itoa(port)
	for connectionRetryCount < maxConnectionRetries {
		conn, err = dialer.Dial("tcp", dest)
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				if e.Temporary() {
					connectionRetryCount++
					time.Sleep(connectionRetrySleep)
					continue
				}
				log.Fatalln("Could not connect to the server:", err)
			default:
				log.Fatalln("Could not connect to the server:", err)
			}
		}
		break
	}
	return conn
}

func initCollectRequest(conn net.Conn) (CommandResponse, error) {
	var commandResponse CommandResponse
	commandRequest := CommandRequest{
		Command: "collect",
	}
	if err := json.NewEncoder(conn).Encode(&commandRequest); err != nil {
		return commandResponse, err
	}
	err := json.NewDecoder(conn).Decode(&commandResponse)
	if err != nil {
		return commandResponse, err
	}
	return commandResponse, nil
}

func sendLog(conn net.Conn, notification Notification) (CollectCmdResponse, error) {
	var collectResp CollectCmdResponse
	collectCmdPayload := CollectCmdPayload{
		Record:      notification.LogEvent,
		Timestamp:   notification.Timestamp,
		Tags:        notification.Tags,
		Destination: notification.Destination,
	}
	if err := json.NewEncoder(conn).Encode(&collectCmdPayload); err != nil {
		return collectResp, err
	}
	err := json.NewDecoder(conn).Decode(&collectResp)
	if err != nil {
		return collectResp, err
	}
	return collectResp, nil
}
