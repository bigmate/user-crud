package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"user-crud/internal/config"
	"user-crud/internal/repository/kafka"
	"user-crud/internal/repository/postgres"
	"user-crud/internal/transport/grpc"
	"user-crud/pkg/app"
	"user-crud/pkg/closer"
	"user-crud/pkg/healthcheck"
	"user-crud/pkg/interceptors"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to init config: %v\n", err)
	}

	postgresClient, err := postgres.NewClient(ctx, conf)
	if err != nil {
		log.Fatalf("failed to init postgresql client: %v", err)
	}

	kafkaClient, err := kafka.NewClient(conf)
	if err != nil {
		log.Fatalf("failed to init kafka client: %v", err)
	}

	grpcApp := grpc.NewServer(
		postgresClient,
		kafkaClient,
		grpc.WithPort("8081"),
		grpc.WithHTTP("8080"),
		grpc.WithUnaryInterceptor(interceptors.Validate),
	)
	if err != nil {
		log.Fatalf("failed to start grpc with http: %v\n", err)
	}

	healthChecker := healthcheck.New(ctx,
		healthcheck.WithResource(postgresClient),
		healthcheck.WithResource(kafkaClient),
		healthcheck.WithTimeout(time.Second*5),
	)

	runner := app.NewRunner(grpcApp, healthChecker, closer.NewCloser())

	if err = runner.Run(ctx); err != nil {
		log.Fatalf("failed to run apps: %v", err)
	}
}
