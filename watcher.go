package logcollector

import (
	"log"
	"os"

	"github.com/hpcloud/tail"
)

type Watcher struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}

func Watch(fileName string, notificationChannel chan string) {
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

func getOffset(fileName string) tail.SeekInfo {
	fileInfo, _ := os.Stat(fileName)
	return tail.SeekInfo{
		Offset: fileInfo.Size(),
	}
	// TODO: Add advanced offset management logic to maintain last read state
}
