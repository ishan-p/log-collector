package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/ishan-p/log-collector/internal/schema"
)

func TestReadServerConfigJson(t *testing.T) {
	configFile := "../../../tests/server.config.json"
	expectedConfig := schema.ServerConfig{
		Host:                 "127.0.0.1",
		Port:                 8000,
		ServerWaitTimeSec:    3,
		ConnectionIdleTime:   60 * time.Second,
		MaxConnectionRetries: 5,
		SleepRetryDuration:   10 * time.Millisecond,
		Storage: schema.StorageConfig{
			Filesystem: schema.FsStorageConfig{
				BaseDir: "/tmp/collector/logs",
			},
		},
	}

	config := ReadJSON(configFile)
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Fatalf("Server config JSON unmarshal failed!\n Got: \n%v\n Want: \n%v\n", config, expectedConfig)
	}
}
