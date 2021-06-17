package logcollector

import "fmt"

func RunAgent(configFile string) {
	config := readClientConfigJSON(configFile)
	done := make(chan bool)
	for _, watcher := range config.Watchers {
		fmt.Printf("Watching file %s\n", watcher.FileName)
		go Watch(watcher)
	}
	<-done
}
