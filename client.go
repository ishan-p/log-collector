package logcollector

import (
	"log"
)

func RunAgent(configFile string) {
	config := readClientConfigJSON(configFile)
	done := make(chan bool)
	for _, watcher := range config.Watchers {
		log.Printf("Watching file %s\n", watcher.FileName)
		go Watch(watcher)
	}
	<-done
}
