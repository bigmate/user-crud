package notifier

import (
	"io"
)

// Service is an upper layer to interact with other services
type Service interface {
	Publish(topic string, reader io.Reader) error
}
