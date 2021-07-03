package server

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

type S3StorageConfig struct {
	FirehosStream string `json:"firehose_stream"`
	AWSRegion     string `json:"aws_region"`
}

type S3Store struct {
	S3StorageConfig
}

func NewS3Store(firehoseStream string, awsRegion string) S3Store {
	s3Config := S3StorageConfig{
		FirehosStream: firehoseStream,
		AWSRegion:     awsRegion,
	}
	s3store := S3Store{
		s3Config,
	}
	return s3store
}

func (s3 S3Store) write(data []byte) (bool, error) {
	sess := session.Must(session.NewSession())
	firehoseService := firehose.New(sess, aws.NewConfig().WithRegion(s3.AWSRegion))

	recordInput := &firehose.PutRecordInput{}
	recordInput = recordInput.SetDeliveryStreamName(s3.FirehosStream)

	record := &firehose.Record{Data: data}
	recordInput = recordInput.SetRecord(record)

	_, err := firehoseService.PutRecord(recordInput)
	if err != nil {
		log.Println("Put to firehose failed with error: ", err)
		return false, err
	}
	return true, nil
}
