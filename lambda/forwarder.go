package lambda

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/config"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/connector"
	"github.com/jacob-elektronik/rabbit-amazon-forwarder/forwarder"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

const (
	// Type forwarder type
	Type = "Lambda"
)

// Forwarder forwarding client
type Forwarder struct {
	name         string
	lambdaClient lambdaiface.LambdaAPI
	function     string
}

// CreateForwarder creates instance of forwarder
func CreateForwarder(entry config.AmazonEntry, lambdaClient ...lambdaiface.LambdaAPI) forwarder.Client {
	var client lambdaiface.LambdaAPI
	if len(lambdaClient) > 0 {
		client = lambdaClient[0]
	} else {
		client = lambda.New(connector.CreateAWSSession())
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

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	params := &lambda.InvokeInput{
		FunctionName: aws.String(f.function),
		Payload:      []byte(message),
	}
	resp, err := f.lambdaClient.InvokeWithContext(ctx, params)
	if err != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"error":         err.Error()}).Error("Could not forward message")
		return err
	}
	if resp.FunctionError != nil {
		log.WithFields(log.Fields{
			"forwarderName": f.Name(),
			"functionError": *resp.FunctionError}).Errorf("Could not forward message")
		return err
	}
	log.WithFields(log.Fields{
		"forwarderName": f.Name(),
		"statusCode":    resp.StatusCode}).Info("Forward succeeded")
	return nil
}
