package sqs

import (
	"errors"

	"github.com/jacob-elektronik/rabbit-amazon-forwarder/connector"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/fdegner/go-spanctx"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/config"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/forwarder"
	"github.com/opentracing/opentracing-go"
)

const (
	// Type forwarder type
	Type = "SQS"
)

// Forwarder forwarding client
type Forwarder struct {
	name      string
	sqsClient sqsiface.SQSAPI
	queue     string
}

// CreateForwarder creates instance of forwarder
func CreateForwarder(entry config.AmazonEntry, sqsClient ...sqsiface.SQSAPI) forwarder.Client {
	var client sqsiface.SQSAPI
	if len(sqsClient) > 0 {
		client = sqsClient[0]
	} else {
		client = sqs.New(connector.CreateAWSSession())
	}
	f := Forwarder{entry.Name, client, entry.Target}
	log.WithField("forwarderName", f.Name()).Info("Created forwarder")
	return f
}

// Name forwarder name
func (f Forwarder) Name() string {
	return f.name
}

// Push pushes message to forwarding infrastructure
func (f Forwarder) Push(span opentracing.Span, message string) error {
	if message == "" {
		err := errors.New(forwarder.EmptyMessageError)
		return err
	}
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(message), // Required
		QueueUrl:    aws.String(f.queue), // Required
	}
	err := spanctx.AddToSQSMessageInput(span.Context(), params)
	if err != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"error":         err.Error()}).Error("Could not inject span context into SQS message attributes")
	}

	resp, err := f.sqsClient.SendMessage(params)
	if err != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"error":         err.Error()}).Error("Could not forward message")
		return err
	}
	log.WithFields(log.Fields{
		"forwarderName": f.Name(),
		"responseID":    resp.MessageId}).Info("Forward succeeded")
	return nil
}
