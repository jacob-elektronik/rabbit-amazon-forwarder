package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/opentracing/opentracing-go"
)

type sqsAttributeCarrier map[string]*sqs.MessageAttributeValue

func (c sqsAttributeCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, val := range c {
		if *val.DataType != "String" {
			continue
		}
		if err := handler(k, *val.StringValue); err != nil {
			return err
		}
	}
	return nil
}

func (c sqsAttributeCarrier) Set(key, val string) {
	c[key].DataType = aws.String("String")
	c[key].StringValue = aws.String(val)
}

func injectSpanContext(span opentracing.Span, pubInput *sqs.SendMessageInput) error {
	if span == nil {
		return nil
	}
	c := sqsAttributeCarrier(pubInput.MessageAttributes)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, c)
}

func extractSpanContext(pubInput *sqs.SendMessageInput) (opentracing.SpanContext, error) {
	c := sqsAttributeCarrier(pubInput.MessageAttributes)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, c)
}
