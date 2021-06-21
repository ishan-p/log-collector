package logcollector

import (
	"log"
)

func RunAgent(configFile string) {
	config := readClientConfigJSON(configFile)
	done := make(chan bool)
	for _, watcher := range config.Watchers {
		log.Printf("Watching file %s\n", watcher.FileName)
		notificationChannel := make(chan string)
		defer close(notificationChannel)
		go Watch(watcher.FileName, notificationChannel)
		go Notify(notificationChannel)
	}
	<-done
}

func Notify(ch chan string) {
	for notification := range ch {
		log.Println(notification)
	}
}
