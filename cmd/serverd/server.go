package main

import (
	"github.com/ishan-p/log-collector/internal/server"
)

func main() {
	server.Run("./server.config.json")
}
