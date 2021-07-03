package schema

import "time"

type ServerConfig struct {
	Host                 string        `json:"host"`
	Port                 int           `json:"port"`
	ServerWaitTimeSec    int           `json:"wait_time"`
	ConnectionIdleTime   time.Duration `json:"connection_idle_time"`
	MaxConnectionRetries int           `json:"max_connection_retries"`
	SleepRetryDuration   time.Duration `json:"sleep_retry"`
	Storage              StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Filesystem FsStorageConfig `json:"filesystem"`
	S3         S3StorageConfig `json:"s3"`
}

type S3StorageConfig struct {
	FirehosStream string `json:"firehose_stream"`
	AWSRegion     string `json:"aws_region"`
}

type FsStorageConfig struct {
	BaseDir string `json:"base_dir"`
}
