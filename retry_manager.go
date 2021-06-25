package logcollector

import (
	"log"
	"sync"
	"time"
)

type RetryConfiguration struct {
	MaxRetries   int `json:"max_retries"`
	MaxQueueSize int `json:"max_queue_size"`
}

func defaultRetryConfiguration() RetryConfiguration {
	var retryConfig RetryConfiguration
	retryConfig.MaxQueueSize = 500
	retryConfig.MaxRetries = 3
	return retryConfig
}

type RetryManager struct {
	MaxRetries   int
	MaxQueueSize int
	RetryQueue   []Notification
	mu           sync.Mutex
}

func (manager *RetryManager) retry(retryCh chan Notification, server ServerConnection) {
	for {
		var currentRetryItem Notification
		manager.mu.Lock()
		queueSize := len(manager.RetryQueue)
		if queueSize > 0 {
			currentRetryItem = manager.RetryQueue[0]
		}
		manager.mu.Unlock()
		if queueSize > 0 {
			backOffTime := 10 * time.Duration(currentRetryItem.RetryAttempt) * time.Second
			if currentRetryItem.LastRetry <= time.Now().Add(-backOffTime).Unix() {
				currentRetryItem.LastRetry = time.Now().Unix()
				log.Printf("Slept for %v - Retry attempt %v - %v\n", backOffTime, currentRetryItem.RetryAttempt, currentRetryItem.LogEvent)
				sendToServer(currentRetryItem, server, retryCh)
				manager.mu.Lock()
				manager.RetryQueue = manager.RetryQueue[1:]
				manager.mu.Unlock()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func retryMangager(retryCh chan Notification, server ServerConnection, retryConfig RetryConfiguration) {
	manager := RetryManager{
		MaxRetries:   retryConfig.MaxRetries,
		MaxQueueSize: retryConfig.MaxQueueSize,
		RetryQueue:   make([]Notification, 0, retryConfig.MaxQueueSize),
	}
	go manager.retry(retryCh, server)
	for notification := range retryCh {
		if manager.MaxRetries > notification.RetryAttempt {
			updatedNotification := Notification{
				LogEvent:     notification.LogEvent,
				Timestamp:    notification.Timestamp,
				Tags:         notification.Tags,
				Destination:  notification.Destination,
				RetryAttempt: notification.RetryAttempt + 1,
				LastRetry:    notification.LastRetry,
			}
			manager.mu.Lock()
			if cap(manager.RetryQueue) > len(manager.RetryQueue) {
				manager.RetryQueue = append(manager.RetryQueue, updatedNotification)
			} else {
				log.Println("Retry queue full")
			}
			manager.mu.Unlock()
		}
	}
}
