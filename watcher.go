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

func Watch(watcherSource Watcher) {
	seekInfo := getOffset(watcherSource.FileName)
	t, err := tail.TailFile(watcherSource.FileName, tail.Config{
		Follow:   true,
		Location: &seekInfo,
	})
	if err != nil {
		log.Println("Could not tail file - ", watcherSource.FileName)
	}
	for line := range t.Lines {
		log.Println(line.Text)
	}
}

func getOffset(fileName string) tail.SeekInfo {
	fileInfo, _ := os.Stat(fileName)
	return tail.SeekInfo{
		Offset: fileInfo.Size(),
	}
}
