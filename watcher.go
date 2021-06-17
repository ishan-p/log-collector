package logcollector

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	FileName    string              `json:"watch"`
	Tags        []map[string]string `json:"tags"`
	Destination string              `json:"destination"`
}

type eventHandler struct {
	writeCh chan string
}

func Watch(watcherSource Watcher) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("ERROR", err)
	}
	defer watcher.Close()

	err = watcher.Add(watcherSource.FileName)
	if err != nil {
		log.Println("ERROR", err)
	}

	eventChan := make(chan fsnotify.Event)
	defer close(eventChan)
	go handleEvents(eventChan, watcherSource)

	for {
		select {
		case event, _ := <-watcher.Events:
			go func() { eventChan <- event }()

		case err, _ := <-watcher.Errors:
			log.Println("ERROR", err)
		}
	}
}

func handleEvents(eventCh chan fsnotify.Event, watcherSource Watcher) {
	handler := eventHandler{
		writeCh: make(chan string),
	}
	runEventActions(handler, watcherSource)
	for event := range eventCh {
		if event.Op&fsnotify.Write == fsnotify.Write {
			handler.writeCh <- event.Name
		}
		// TODO: Handle other notfification types
	}
}

func runEventActions(handler eventHandler, watcherSource Watcher) {
	go actionOnWrite(handler.writeCh, watcherSource)
	// TODO: Handle other event actions
}

func actionOnWrite(eventName chan string, watcherSource Watcher) {
	var line string
	fileInfo, _ := os.Stat(watcherSource.FileName)
	currentOffset := fileInfo.Size()
	for e := range eventName {
		line, currentOffset = readLine(e, currentOffset)
		log.Println(line)
		// Send log to the server
	}
}

func readLine(filename string, offset int64) (string, int64) {
	var part []byte
	var prefix bool
	var line string
	fd, err := os.Open(filename)
	defer fd.Close()
	if err != nil {
		log.Fatalln("Failed to open the file")
	}
	_, err = fd.Seek(offset, 0)
	if err != nil {
		log.Fatalln("Failed to read file at given offset")
	}

	reader := bufio.NewReader(fd)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			line = buffer.String()
			offset = offset + int64(buffer.Len()) + 1
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return line, offset
}
