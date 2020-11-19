package rabbitmq

import (
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

type amqpHeaderCarrier map[string]interface{}

func (c amqpHeaderCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, val := range c {
		v, ok := val.(string)
		if !ok {
			continue
		}
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (c amqpHeaderCarrier) Set(key, val string) {
	c[key] = val
}

func injectSpanContext(span opentracing.Span, table amqp.Table) error {
	c := amqpHeaderCarrier(table)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, c)
}

func extractSpanContext(table amqp.Table) (opentracing.SpanContext, error) {
	c := amqpHeaderCarrier(table)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, c)
}
