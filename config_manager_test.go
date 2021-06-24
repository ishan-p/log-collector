package logcollector

import (
	"reflect"
	"testing"
	"time"
)

func TestReadClientConfigJson(t *testing.T) {
	configFile := "tests/client.config.json"
	var expectedConfig ClientConfig
	expectedConfig.Watchers = []Watcher{
		{
			FileName: "/var/log/nginx.log",
			Tags: []map[string]string{
				{
					"key":   "Type",
					"value": "Nginx",
				},
			},
			Destination: "filesystem",
		},
	}
	expectedConfig.Collector = ServerConnection{
		Host:                 "127.0.0.1",
		Port:                 8000,
		ServerWaitTimeSec:    3,
		ConnectionIdleTime:   45 * time.Second,
		MaxConnectionRetries: 5,
		SleepRetryDuration:   10 * time.Millisecond,
	}
	expectedConfig.RetryParams = RetryConfiguration{
		MaxRetries:   5,
		MaxQueueSize: 500,
	}

	config := readClientConfigJSON(configFile)
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Client config JSON unmarshal failed!\n Got: \n%v\n Want: \n%v\n", config, expectedConfig)
	}
}

func TestReadServerConfigJson(t *testing.T) {
	configFile := "tests/server.config.json"
	serverConnection := ServerConnection{
		Host:                 "127.0.0.1",
		Port:                 8000,
		ServerWaitTimeSec:    3,
		ConnectionIdleTime:   60 * time.Second,
		MaxConnectionRetries: 5,
		SleepRetryDuration:   10 * time.Millisecond,
	}
	expectedConfig := ServerConfig{
		serverConnection,
		StorageConfig{
			Filesystem: FsStorageConfig{
				BaseDir: "/tmp/collector/logs",
			},
		},
	}

	config := readServerConfigJSON(configFile)
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Server config JSON unmarshal failed!\n Got: \n%v\n Want: \n%v\n", config, expectedConfig)
	}
}
