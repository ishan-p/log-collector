package main

import (
	"github.com/ishan-p/log-collector/internal/client"
)

func main() {
	client.Run("./client.config.json")
}
