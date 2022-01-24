package healthcheck

import (
	"time"
)

type Option func(hc *healthCheck)

func WithPort(port string) Option {
	return func(hc *healthCheck) {
		hc.port = port
	}
}

func WithPath(path string) Option {
	return func(hc *healthCheck) {
		hc.path = path
	}
}

func WithResource(resource Resource) Option {
	return func(hc *healthCheck) {
		hc.resources = append(hc.resources, resource)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(hc *healthCheck) {
		hc.timeout = timeout
	}
}

func defaultHealthCheck() *healthCheck {
	return &healthCheck{
		port:    "8082",
		path:    "/health",
		timeout: time.Second * 10,
	}
}
