package client

import (
	"time"

	"github.com/ishan-p/log-collector/internal/schema"
)

type Notification struct {
	LogEvent     string
	Timestamp    int64
	Tags         []map[string]string
	Destination  string
	RetryAttempt int
	LastRetry    int64
}

func notify(ch chan string, server schema.CollectorConfig, destination string, tags []map[string]string, retryChannel chan Notification) {
	for event := range ch {
		notification := Notification{
			LogEvent:     event,
			Timestamp:    time.Now().Unix(),
			Tags:         tags,
			Destination:  destination,
			RetryAttempt: 0,
			LastRetry:    time.Now().Unix(),
		}
		sendToServer(notification, server, retryChannel)
	}
}
