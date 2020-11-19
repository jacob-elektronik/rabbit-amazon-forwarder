package sns

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/config"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/connector"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/forwarder"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

const (
	// Type forwarder type
	Type = "SNS"
)

// Forwarder forwarding client
type Forwarder struct {
	name      string
	snsClient snsiface.SNSAPI
	topic     string
}

// CreateForwarder creates instance of forwarder
func CreateForwarder(entry config.AmazonEntry, snsClient ...snsiface.SNSAPI) forwarder.Client {
	var client snsiface.SNSAPI
	if len(snsClient) > 0 {
		client = snsClient[0]
	} else {
		client = sns.New(connector.CreateAWSSession())
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
		span.SetTag("error", err.Error())
		return err
	}
	params := &sns.PublishInput{
		Message:   aws.String(message),
		TargetArn: aws.String(f.topic),
	}
	err := injectSpanContext(span, params)
	if err != nil {
		span.SetTag("error", err.Error())
		return err
	}

	resp, err := f.snsClient.Publish(params)
	if err != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"error":         err.Error()}).Error("Could not forward message")
		span.SetTag("error", err.Error())
		return err
	}
	log.WithFields(log.Fields{
		"forwarderName": f.Name(),
		"responseID":    resp.MessageId}).Info("Forward succeeded")
	return nil
}
