package logcollector

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestWatcherWriteReadCount(t *testing.T) {
	logEntries := 1000
	watchTimeout := time.Duration(5)
	logFile := "tests/sample.log"
	notificationChannel := make(chan string)
	go watch(logFile, notificationChannel)

	time.Sleep(1 * time.Second)
	go func() {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		for i := 0; i < logEntries; i++ {
			if _, err = f.WriteString(fmt.Sprintf("Test %s\n", fmt.Sprint(i))); err != nil {
				panic(err)
			}
		}
		f.Close()
	}()

	func() {
		var receivedLogEvents int
		for {

			select {
			case <-notificationChannel:
				receivedLogEvents++
				if receivedLogEvents == logEntries {
					return
				}
			case <-time.After(watchTimeout * time.Second):
				log.Fatalln("Watcher add-read count test failed - Timeout.")
			}
		}
	}()
}
