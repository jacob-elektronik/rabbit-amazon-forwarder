package lambda

import (
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

// XXX: This code firmly believes it owns the Lambda's InvokeInput.ClientContext

func injectSpanContext(span opentracing.Span, input *lambda.InvokeInput) error {
	if span == nil {
		return nil
	}
	ctx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		// opentracing.Span doesn't require String(), but we do
		return errors.New("span context implementation not supported")
	}
	if *input.InvocationType != lambda.InvocationTypeRequestResponse {
		// Setting InvokeInput.ClientContext is only supported for RequestResponse type
		return errors.New("this span injection only works for invocation type RequestResponse")
	}
	var buf []byte
	base64.StdEncoding.Encode(buf, []byte(ctx.String()))
	input.ClientContext = aws.String(string(buf))
	return nil
}

func extractSpanContext(input *lambda.InvokeInput) (opentracing.SpanContext, error) {
	if input == nil {
		return nil, nil
	}
	buf := []byte(*input.ClientContext)
	var rawSpan []byte
	base64.StdEncoding.Decode(buf, rawSpan)
	ctx, err := jaeger.ContextFromString(string(rawSpan))
	if err != nil {
		return nil, err
	}
	return opentracing.SpanContext(ctx), nil
}
