package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"users/gen/pb"
	"users/postgres"
	"users/services/users"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Params struct {
	GrpcPort int
	HttpPort int
	Postgres postgres.IPostgres
}

func Run(params Params) {
	s := grpc.NewServer()

	users.Register(s, params.Postgres)

	go runGrpc(params.GrpcPort, s)
	runHttp(params.HttpPort, params.GrpcPort)
}

func runGrpc(port int, s *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func runHttp(httpport int, grpcport int) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	grpcadress := fmt.Sprintf("localhost:%d", grpcport)
	err := pb.RegisterUsersHandlerFromEndpoint(ctx, mux, grpcadress, opts)
	if err != nil {
		panic(err)
	}

	log.Printf("http server listening at %d", httpport)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpport), mux); err != nil {
		panic(err)
	}
}
