package tcp

import (
	"log"
	"net"
	"strconv"
	"time"
)

type Handler interface {
	Handle(conn net.Conn)
}

type Server struct {
	Type                 string
	Port                 int
	Listener             net.Listener
	SleepBetweenRetries  time.Duration
	MaxConnectionRetries int
	ConnectionIdleTime   time.Duration
}

func NewServer(port int) Server {
	// TODO: Accept make other params configurable
	server := Server{
		Type:                 "tcp",
		Port:                 port,
		SleepBetweenRetries:  10 * time.Millisecond,
		MaxConnectionRetries: 5,
		ConnectionIdleTime:   45 * time.Second,
	}
	return server
}

func (server Server) Listen() (net.Listener, error) {
	src := ":" + strconv.Itoa(server.Port)
	listener, err := net.Listen(server.Type, src)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (server Server) Stop() {
	server.Listener.Close()
}

func (server Server) WaitForConnection() (net.Conn, error) {
	var connectionRetryCount int
	sleepBetweenRetry := server.SleepBetweenRetries
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				if e.Temporary() {
					connectionRetryCount++
					if connectionRetryCount > server.MaxConnectionRetries {
						log.Printf("Unable to accept connections after %d retries: %v\n", connectionRetryCount, err)
						return nil, err
					}
					sleepBetweenRetry *= 2
					time.Sleep(sleepBetweenRetry)
					continue
				} else {
					log.Println(err)
					return nil, err
				}
			default:
				conn.Close()
				return nil, err
			}
		}
		if err := conn.SetDeadline(time.Now().Add(server.ConnectionIdleTime)); err != nil {
			log.Println("Failed to set deadline:", err)
			return nil, err
		}
		return conn, nil
	}
}

func (server Server) HandleConnection(conn net.Conn, handler Handler) {
	go handler.Handle(conn)
}
