package forwarder

import "github.com/opentracing/opentracing-go"

const (
	// EmptyMessageError empty error message
	EmptyMessageError = "message is empty"
)

// Client interface to forwarding messages
type Client interface {
	Name() string
	Push(span opentracing.Span, message string) error
}
