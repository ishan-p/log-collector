package logcollector

import (
	"reflect"
	"testing"
)

func TestReadClientConfigJson(t *testing.T) {
	configFile := "tests/sample.config.json"
	var expectedConfig Config
	expectedConfig.Watchers = []Watcher{
		{
			FileName: "/var/log/nginx.log",
			Tags: []map[string]string{
				{
					"key":   "Type",
					"value": "Nginx",
				},
			},
			Destination: "S3",
		},
	}
	expectedConfig.Collector = ServerConnection{
		Host:              "127.0.0.1",
		Port:              8000,
		ServerWaitTimeSec: 3,
	}
	expectedConfig.RetryParams = RetryConfiguration{
		MaxRetries:   5,
		MaxQueueSize: 500,
	}

	config := readClientConfigJSON(configFile)
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatal("Client config JSON unmarshal failed!")
	}
}
