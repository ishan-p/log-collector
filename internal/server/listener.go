package server

import (
	"log"

	serverConfig "github.com/ishan-p/log-collector/internal/config/server"
	"github.com/ishan-p/log-collector/internal/platform/tcp"
)

func Run(configFile string) {
	var err error
	config := serverConfig.ReadJSON(configFile)
	server := tcp.NewServer(config.Port)
	server.Listener, err = server.Listen()
	if err != nil {
		log.Fatalln(err)
	}
	handler := RequestHandler{
		Config: config,
	}
	for {
		conn, err := server.WaitForConnection()
		if err != nil {
			log.Fatalln(err)
		}
		server.HandleConnection(conn, handler)
	}
}
