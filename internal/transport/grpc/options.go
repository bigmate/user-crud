package grpc

import (
	"google.golang.org/grpc"
)

type parameters struct {
	grpcPort          string
	exposeHttp        bool
	httpPort          string
	unaryInterceptors []grpc.UnaryServerInterceptor
}

func defaultParameters() *parameters {
	return &parameters{
		grpcPort:   "8081",
		exposeHttp: false,
		httpPort:   "8080",
	}
}

type Option func(p *parameters)

func WithPort(port string) Option {
	return func(p *parameters) {
		p.grpcPort = port
	}
}

func WithHTTP(port string) Option {
	return func(p *parameters) {
		p.exposeHttp = true
		p.httpPort = port
	}
}

func WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) Option {
	return func(p *parameters) {
		p.unaryInterceptors = append(p.unaryInterceptors, interceptor)
	}
}
