package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grpc_sqlc/internal/api/handlers"
	mainapi "github.com/grpc_sqlc/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	s := grpc.NewServer()
	mainapi.RegisterUserServiceServer(s, &handlers.Server{})

	fmt.Println("gRPC Server is running on port", ":50052")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Printf("Error listen to port %v", err)
		return
	}

	// not needed for grpcurl to work, but can be useful for debugging with grpc_client
	reflection.Register(s)

	// listen to configuration
	err = s.Serve(lis)
	if err != nil {
		log.Printf("Error serve %s", err)
		return
	}
}
