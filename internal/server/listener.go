package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	serverConfig "github.com/ishan-p/log-collector/internal/config/server"
	"github.com/ishan-p/log-collector/internal/schema"
)

func Run(configFile string) {
	config := serverConfig.ReadJSON(configFile)
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
		go handleConnection(conn, config.Storage)
	}
}

func handleConnection(conn net.Conn, storageConfig schema.StorageConfig) {
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
		err = handleRequest(command.Command, conn, storageConfig)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func waitForCommand(conn net.Conn) (schema.CommandRequest, error) {
	decoder := json.NewDecoder(conn)
	var request schema.CommandRequest
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}

func beginTransaction(command string, conn net.Conn) error {
	initRequest := schema.CommandResponse{
		Command: command,
		Begin:   true,
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&initRequest); err != nil {
		return err
	}
	return nil
}

func handleRequest(activeCommand string, conn net.Conn, storageConfig schema.StorageConfig) error {
	switch activeCommand {
	case "collect":
		err := handleCollect(conn, storageConfig)
		return err
	default:
		err := fmt.Errorf("Invalid command")
		return err
	}
}
