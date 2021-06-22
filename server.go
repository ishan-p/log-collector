package logcollector

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
)

type ServerConnection struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	ServerWaitTimeSec int    `json:"wait_time"`
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
	return connection
}

func RunServer() {
	port := 8888
	src := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on ", src)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error: ", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	for {
		command, err := waitForCommand(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Closing connection")
				break
			}
			log.Println(err)
		}
		beginTransaction(command.Command, conn)
		handleRequest(command.Command, conn)
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

func beginTransaction(command string, conn net.Conn) {
	initRequest := CommandResponse{
		Command: command,
		Begin:   true,
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&initRequest); err != nil {
		log.Println(err)
	}
}

func handleRequest(activeCommand string, conn net.Conn) {
	switch activeCommand {
	case "collect":
		handleCollect(conn)
	default:
		log.Println("Unrecognized command")
	}
}

func handleCollect(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	var collectPayload CollectCmdPayload
	if err := decoder.Decode(&collectPayload); err != nil {
		log.Println(err)
	}
	ack := collect(collectPayload)
	collectResponse := CollectCmdResponse{Ack: ack}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&collectResponse); err != nil {
		log.Println(err)
	}
}

func collect(record CollectCmdPayload) bool {
	log.Println(record)
	return true
}
