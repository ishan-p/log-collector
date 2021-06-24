package logcollector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type ServerConnection struct {
	Host                 string        `json:"host"`
	Port                 int           `json:"port"`
	ServerWaitTimeSec    int           `json:"wait_time"`
	ConnectionIdleTime   time.Duration `json:"connection_idle_time"`
	MaxConnectionRetries int           `json:"max_connection_retries"`
	SleepRetryDuration   time.Duration `json:"sleep_retry"`
}

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

func defaultServerConnection() ServerConnection {
	var connection ServerConnection
	connection.Host = "127.0.0.1"
	connection.Port = 8888
	connection.ServerWaitTimeSec = 3
	connection.ConnectionIdleTime = 45 * time.Second
	connection.MaxConnectionRetries = 5
	connection.SleepRetryDuration = 10 * time.Millisecond
	return connection
}

func RunServer(configFile string) {
	config := readServerConfigJSON(configFile)
	src := ":" + strconv.Itoa(config.Port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on ", src)
	defer listener.Close()

	var connectionRetryCount int
	connectionRetrySleep := config.SleepRetryDuration
	for {
		conn, err := listener.Accept()
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				if e.Temporary() {
					connectionRetryCount++
					if connectionRetryCount > config.MaxConnectionRetries {
						log.Printf("Unable to accept connections after %d retries: %v\n", connectionRetryCount, err)
						return
					}
					connectionRetrySleep *= 2
					time.Sleep(connectionRetrySleep)
				} else {
					log.Fatalln(err)
				}
			default:
				conn.Close()
				log.Fatalln(err)
			}
			connectionRetryCount = 0
			connectionRetrySleep = config.SleepRetryDuration
		}
		if err := conn.SetDeadline(time.Now().Add(config.ConnectionIdleTime)); err != nil {
			log.Println("Failed to set deadline:", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()
	for {
		command, err := waitForCommand(conn)
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				if e.Timeout() {
					log.Println("Timeout - Did not recieve a request from the client. Closing connection.")
				}
				log.Println(err)
				return
			default:
				if err == io.EOF {
					log.Println("Closing connection on client's request")
					return
				} else {
					log.Println(err)
					continue
				}
			}
		}
		err = beginTransaction(command.Command, conn)
		if err != nil {
			log.Println("Failed to initiate transaction: ", err)
			return
		}
		err = handleRequest(command.Command, conn)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func waitForCommand(conn net.Conn) (CommandRequest, error) {
	decoder := json.NewDecoder(conn)
	var request CommandRequest
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}

func beginTransaction(command string, conn net.Conn) error {
	initRequest := CommandResponse{
		Command: command,
		Begin:   true,
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&initRequest); err != nil {
		return err
	}
	return nil
}

func handleRequest(activeCommand string, conn net.Conn) error {
	switch activeCommand {
	case "collect":
		err := handleCollect(conn)
		return err
	default:
		err := fmt.Errorf("Invalid command")
		return err
	}
}

func handleCollect(conn net.Conn) error {
	decoder := json.NewDecoder(conn)
	var collectPayload CollectCmdPayload
	if err := decoder.Decode(&collectPayload); err != nil {
		log.Println("Failed to read the request")
		return err
	}
	ack := collect(collectPayload)
	collectResponse := CollectCmdResponse{Ack: ack}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&collectResponse); err != nil {
		log.Println("Failed to send the ack")
		return err
	}
	return nil
}

func collect(record CollectCmdPayload) bool {
	log.Println(record)
	return true
}
