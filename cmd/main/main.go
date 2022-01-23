package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"user-crud/internal/config"
	"user-crud/internal/transport/grpc"
	"user-crud/pkg/app"
	"user-crud/pkg/closer"
	"user-crud/pkg/interceptors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v\n", err)
	}

	grpcApp, err := grpc.NewServer(ctx, conf,
		grpc.WithHTTP("8080"),
		grpc.WithUnaryInterceptor(interceptors.Validate),
	)
	if err != nil {
		log.Fatalf("failed to start grpc with http: %v\n", err)
	}

	runner := app.NewRunner(grpcApp, closer.NewCloser())

	if err = runner.Run(ctx); err != nil {
		log.Fatalf("failed to run apps: %v", err)
	}
}
