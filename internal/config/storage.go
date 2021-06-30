package config

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
