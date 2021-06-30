package client

import (
	"log"
	"net"
	"testing"
	"time"

	clientConfig "github.com/ishan-p/log-collector/internal/config/client"
	"github.com/ishan-p/log-collector/internal/server"
)

func init() {
	go server.Run("../../tests/server.config.json")
	time.Sleep(10 * time.Millisecond)
}

func TestServerConnection(t *testing.T) {
	var emptyConn net.Conn
	config := clientConfig.ReadJSON("../../tests/client.config.json")
	conn, err := initServerConnection(config.Collector.Host, config.Collector.Port)
	defer conn.Close()
	if conn == nil || conn == emptyConn || err != nil {
		log.Fatalln("Server connection failed")
	}
}

func TestInitCollectRequest(t *testing.T) {
	config := clientConfig.ReadJSON("../../tests/client.config.json")
	conn, _ := initServerConnection(config.Collector.Host, config.Collector.Port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
}

func TestSendFilesystemLog(t *testing.T) {
	config := clientConfig.ReadJSON("../../tests/client.config.json")
	conn, _ := initServerConnection(config.Collector.Host, config.Collector.Port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
	notification := Notification{
		LogEvent:    "test log",
		Destination: "filesystem",
		Timestamp:   time.Now().Unix(),
	}
	collectAck, err := sendLog(conn, notification)
	if err != nil {
		log.Fatalln("Failed to send log")
	}
	if !collectAck.Ack {
		log.Fatalln("Did not receive collect log Ack")
	}
}
