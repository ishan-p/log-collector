package client

import (
	"log"

	clientConfig "github.com/ishan-p/log-collector/internal/config/client"
)

func Run(configFile string) {
	config := clientConfig.ReadJSON(configFile)
	done := make(chan bool)
	retryChannel := make(chan Notification)
	go retryMangager(retryChannel, config.Collector, config.RetryParams)
	defer close(retryChannel)
	for _, watcher := range config.Watchers {
		log.Printf("Watching file %s\n", watcher.FileName)
		watcher.Tags = append(watcher.Tags, map[string]string{"source_file": watcher.FileName})
		notificationChannel := make(chan string)
		defer close(notificationChannel)
		go watch(watcher.FileName, notificationChannel)
		go notify(notificationChannel, config.Collector, watcher.Destination, watcher.Tags, retryChannel)
	}
	<-done
}
