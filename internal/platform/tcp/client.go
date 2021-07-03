package tcp

import (
	"encoding/json"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Host                 string
	Port                 int
	ConnectionTimeout    time.Duration
	ConnectionKeepAlive  time.Duration
	MaxConnectionRetries int
	SleepBetweenRetries  time.Duration
}

func NewClient(host string, port int) Client {
	// Make other params configurable
	client := Client{
		Host:                 host,
		Port:                 port,
		ConnectionTimeout:    time.Second * 300,
		ConnectionKeepAlive:  time.Minute * 5,
		MaxConnectionRetries: 5,
		SleepBetweenRetries:  time.Second * 1,
	}
	return client
}

func (client Client) Connect() (net.Conn, error) {
	var connectionRetryCount int
	var conn net.Conn
	var err error
	dialer := &net.Dialer{
		Timeout:   client.ConnectionTimeout,
		KeepAlive: client.ConnectionKeepAlive,
	}
	defaultRetrySleep := client.SleepBetweenRetries
	connectionRetrySleep := defaultRetrySleep
	dest := client.Host + ":" + strconv.Itoa(client.Port)
	for connectionRetryCount < client.MaxConnectionRetries {
		conn, err = dialer.Dial("tcp", dest)
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				if e.Temporary() {
					connectionRetryCount++
					time.Sleep(connectionRetrySleep)
					continue
				}
				return nil, err
			default:
				return nil, err
			}
		}
		break
	}
	return conn, nil
}

func (client Client) Send(conn net.Conn, request interface{}) (interface{}, error) {
	if err := json.NewEncoder(conn).Encode(&request); err != nil {
		return nil, err
	}
	var response interface{}
	err := json.NewDecoder(conn).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
