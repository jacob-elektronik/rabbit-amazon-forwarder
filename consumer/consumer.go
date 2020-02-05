package consumer

import "github.com/jacob-elektronik/rabbit-amazon-forwarder/forwarder"

// Client intarface for consuming messages
type Client interface {
	Name() string
	Start(forwarder.Client, chan bool, chan bool) error
}
