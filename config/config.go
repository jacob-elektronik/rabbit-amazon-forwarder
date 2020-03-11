package config

const (
	// MappingFile mapping file environment variable
	MappingType     = "MAPPING_TYPE"
	MappingFile     = "MAPPING_FILE"
	MappingEndpoint = "MAPPING_ENDPOINT"
	CaCertFile      = "CA_CERT_FILE"
	CertFile        = "CERT_FILE"
	KeyFile         = "KEY_FILE"
	LogFormat       = "LOG_FORMAT"
	MappingTypeApi  = "api"
	MappingTypeFile = "file"
)

// RabbitEntry RabbitMQ mapping entry
type RabbitEntry struct {
	Type           string           `json:"type"`
	Name           string           `json:"name"`
	ConnectionURL  string           `json:"connection"`
	ExchangeName   string           `json:"topic"`
	QueueName      string           `json:"queue"`
	RoutingKey     string           `json:"routing"`
	RoutingKeys    []string         `json:"routingKeys"`
	BindOnly       bool             `json:"bindOnly"`
	ExchangeConfig []ExchangeConfig `json:"exchangeConfig"`
}

type ExchangeConfig struct {
	ExchangeName string   `json:"topic"`
	RoutingKeys  []string `json:"routingKeys"`
}

// AmazonEntry SQS/SNS mapping entry
type AmazonEntry struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Target string `json:"target"`
}
