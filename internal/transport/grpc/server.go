package grpc

import (
	"context"
	"log"
	"net"
	"net/http"

	"user-crud/internal/repository"
	"user-crud/internal/services/notifier"
	"user-crud/internal/services/user"
	"user-crud/pkg/app"
	usermanager "user-crud/pkg/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	usermanager.UnimplementedUserManagerServer
	param *parameters
	user  user.Service
}

func NewServer(repo repository.User, notify notifier.Service, options ...Option) app.App {
	param := defaultParameters()

	for _, option := range options {
		option(param)
	}

	return &server{
		param: param,
		user:  user.NewService(repo, notify),
	}
}

func (s *server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", "localhost:"+s.param.grpcPort)
	if err != nil {
		return err
	}

	ops := make([]grpc.ServerOption, 0)

	for _, interceptor := range s.param.unaryInterceptors {
		ops = append(ops, grpc.ChainUnaryInterceptor(interceptor))
	}

	srv := grpc.NewServer(ops...)

	usermanager.RegisterUserManagerServer(srv, s)

	if s.param.exposeHttp {
		if err = s.exposeHTTP(ctx); err != nil {
			return err
		}
	}

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	return srv.Serve(listener)
}

func (s *server) exposeHTTP(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := usermanager.RegisterUserManagerHandlerFromEndpoint(ctx, mux, "localhost:"+s.param.grpcPort, opts); err != nil {
		return err
	}

	go func() {
		if err := http.ListenAndServe(":"+s.param.httpPort, mux); err != nil {
			log.Printf("failed to serve http: %v\n", err)
		}
	}()

	return nil
}
