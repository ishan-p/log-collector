package client

import (
	"reflect"
	"testing"

	"github.com/ishan-p/log-collector/internal/schema"
)

func TestReadClientConfigJson(t *testing.T) {
	configFile := "../../../tests/client.config.json"
	var expectedConfig schema.ClientConfig
	expectedConfig.Watchers = []schema.WatcherConfig{
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
	expectedConfig.Collector = schema.CollectorConfig{
		Host:              "127.0.0.1",
		Port:              8000,
		ServerWaitTimeSec: 3,
	}
	expectedConfig.RetryParams = schema.RetryConfig{
		MaxRetries:   5,
		MaxQueueSize: 500,
	}

	config := ReadJSON(configFile)
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Client config JSON unmarshal failed!\n Got: \n%v\n Want: \n%v\n", config, expectedConfig)
	}
}
