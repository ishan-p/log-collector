package server

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/ishan-p/log-collector/internal/schema"
)

type S3StorageConfig struct {
	FirehosStream string `json:"firehose_stream"`
	AWSRegion     string `json:"aws_region"`
}

func writeS3(logEvent schema.CollectCmdPayload, firehoseStream string, awsRegion string) {
	sess := session.Must(session.NewSession())
	firehoseService := firehose.New(sess, aws.NewConfig().WithRegion(awsRegion))

	recordInput := &firehose.PutRecordInput{}
	recordInput = recordInput.SetDeliveryStreamName(firehoseStream)
	jsonLogEvent, err := json.Marshal(logEvent)
	if err != nil {
		log.Println("Unable to encode log as json")
	}
	record := &firehose.Record{Data: jsonLogEvent}
	recordInput = recordInput.SetRecord(record)

	_, err = firehoseService.PutRecord(recordInput)
	if err != nil {
		log.Println("Put to firehose failed with error: ", err)
	}
}
