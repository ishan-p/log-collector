package logcollector

import (
	"log"
	"net"
	"testing"
	"time"
)

func init() {
	go RunServer("./tests/server.config.json")
	time.Sleep(10 * time.Millisecond)
}

func TestServerConnection(t *testing.T) {
	var emptyConn net.Conn
	config := readClientConfigJSON("./tests/client.config.json")
	conn := initServerConnection(config.Collector.Host, config.Collector.Port)
	defer conn.Close()
	if conn == nil || conn == emptyConn {
		log.Fatalln("Server connection failed")
	}
}

func TestInitCollectRequest(t *testing.T) {
	config := readClientConfigJSON("./tests/client.config.json")
	conn := initServerConnection(config.Collector.Host, config.Collector.Port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
}

func TestSendFilesystemLog(t *testing.T) {
	config := readClientConfigJSON("./tests/client.config.json")
	conn := initServerConnection(config.Collector.Host, config.Collector.Port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
	collectAck, err := sendLog(conn, "test log", "filesystem", make([]map[string]string, 0))
	if err != nil {
		log.Fatalln("Failed to send log")
	}
	if !collectAck.Ack {
		log.Fatalln("Did not receive collect log Ack")
	}
}
