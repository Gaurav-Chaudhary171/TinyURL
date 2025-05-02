package main

import (
	"log"
	"net"

	"TinyURL_Refactored/internal/service"
	"TinyURL_Refactored/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Create service instance
	tinyURLService := service.NewTinyURLService()

	grpcServer := grpc.NewServer()

	// Register gRPC service
	proto.RegisterTinyURLServiceServer(grpcServer, tinyURLService)

	log.Println("gRPC Server starting on :7777")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
