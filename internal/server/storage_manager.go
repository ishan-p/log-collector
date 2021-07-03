package server

import (
	"errors"

	"github.com/ishan-p/log-collector/internal/schema"
)

type Storer interface {
	write(data []byte) (bool, error)
}

var InvalidDestinationErr error

func init() {
	InvalidDestinationErr = errors.New("Invalid destination")
}

func NewStorer(destination string, storageConfig schema.StorageConfig) (Storer, error) {
	switch destination {
	case "filesystem":
		fs := NewFileStore(storageConfig.Filesystem.BaseDir)
		return fs, nil
	case "s3":
		s3 := NewS3Store(storageConfig.S3.FirehosStream, storageConfig.S3.AWSRegion)
		return s3, nil
	default:
		return nil, InvalidDestinationErr
	}
}
