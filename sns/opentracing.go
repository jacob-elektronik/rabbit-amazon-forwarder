package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/opentracing/opentracing-go"
)

type snsAttributeCarrier map[string]*sns.MessageAttributeValue

func (c snsAttributeCarrier) ForeachKey(handler func(key, val string) error) error {
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

func (c snsAttributeCarrier) Set(key, val string) {
	c[key].DataType = aws.String("String")
	c[key].StringValue = aws.String(val)
}

func injectSpanContext(span opentracing.Span, pubInput *sns.PublishInput) error {
	if span == nil {
		return nil
	}
	c := snsAttributeCarrier(pubInput.MessageAttributes)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, c)
}

func extractSpanContext(pubInput *sns.PublishInput) (opentracing.SpanContext, error) {
	c := snsAttributeCarrier(pubInput.MessageAttributes)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, c)
}
