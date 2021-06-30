package client

import (
	"log"

	"github.com/hpcloud/tail"
)

type Watcher struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}

func watch(fileName string, notificationChannel chan string) {
	seekInfo := getOffset(fileName)
	t, err := tail.TailFile(fileName, tail.Config{
		Follow:   true,
		Location: &seekInfo,
	})
	if err != nil {
		log.Fatalln("Could not tail file - ", fileName)
	}
	for line := range t.Lines {
		notificationChannel <- line.Text
	}
}
