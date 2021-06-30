package server

import (
	"github.com/ishan-p/log-collector/internal/schema"
)

func collect(record schema.CollectCmdPayload, storageConfig schema.StorageConfig) bool {
	if record.Destination == "filesystem" {
		writeFs(record, storageConfig.Filesystem.BaseDir)
		return true
	} else if record.Destination == "s3" {
		writeS3(record, storageConfig.S3.FirehosStream, storageConfig.S3.AWSRegion)
		return true
	}
	return false
}
