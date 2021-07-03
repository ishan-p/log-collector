package server

import (
	"encoding/json"
	"io"
	"log"
	"net"

	"github.com/ishan-p/log-collector/internal/schema"
)

type RequestHandler struct {
	Config schema.ServerConfig
}

func (handler RequestHandler) Handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()
	for {
		cmd, err := waitForCommand(conn)
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
		command, err := NewCommand(cmd)
		status, err := command.Execute(handler)
		err = command.Reply(conn, status)
	}
}

func waitForCommand(conn net.Conn) (interface{}, error) {
	decoder := json.NewDecoder(conn)
	var request interface{}
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}
