package notifier

import (
	"io"

	"user-crud/internal/config"
	"user-crud/internal/services/notifier/kafka"
)

// Service is an upper layer to interact with other services
type Service interface {
	Publish(topic string, reader io.Reader) error
}

// NewService is a constructor of a Service
func NewService(config *config.Config) (Service, error) {
	return kafka.NewClient(config)
}
