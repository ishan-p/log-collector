package logcollector

import (
	"log"
	"net"
	"testing"
	"time"
)

func init() {
	go RunServer()
	time.Sleep(10 * time.Millisecond)
}

func TestServerConnection(t *testing.T) {
	var emptyConn net.Conn
	host := "127.0.0.1"
	port := 8888
	conn := initServerConnection(host, port)
	defer conn.Close()
	if conn == nil || conn == emptyConn {
		log.Fatalln("Server connection failed")
	}
}

func TestInitCollectRequest(t *testing.T) {
	host := "127.0.0.1"
	port := 8888
	conn := initServerConnection(host, port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
}

func TestSendLog(t *testing.T) {
	host := "127.0.0.1"
	port := 8888
	conn := initServerConnection(host, port)
	commandResp, err := initCollectRequest(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if !commandResp.Begin {
		log.Fatalln("Failed to initiate collect command")
	}
	collectAck, err := sendLog(conn, "test log", "S3", make([]map[string]string, 0))
	if err != nil {
		log.Fatalln("Failed to send log")
	}
	if !collectAck.Ack {
		log.Fatalln("Did not receive collect log Ack")
	}
}
