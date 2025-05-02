package main

import (
	"context"
	"log"
	"time"

	"TinyURL_Refactored/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server with retries
	var conn *grpc.ClientConn
	var err error
	for i := 0; i < 5; i++ {
		conn, err = grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
		log.Printf("Failed to connect, retrying in 1 second... (attempt %d/5)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect after 5 attempts: %v", err)
	}
	defer conn.Close()

	client := proto.NewTinyURLServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test RegisterUser
	registerResp, err := client.RegisterUser(ctx, &proto.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "testuser",
	})
	if err != nil {
		log.Printf("RegisterUser failed: %v", err)
	} else {
		log.Printf("RegisterUser response: %v", registerResp)
	}

	// Test Login
	loginResp, err := client.Login(ctx, &proto.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	})
	if err != nil {
		log.Printf("Login failed: %v", err)
	} else {
		log.Printf("Login response: %v", loginResp)
	}

	// Test ShortenURL
	shortenResp, err := client.ShortenURL(ctx, &proto.ShortenURLRequest{
		Url:      "https://example.com",
		Username: "testuser",
	})
	if err != nil {
		log.Printf("ShortenURL failed: %v", err)
	} else {
		log.Printf("ShortenURL response: %v", shortenResp)
	}

	// Test ExtendURL only if ShortenURL was successful
	if shortenResp != nil {
		extendResp, err := client.ExtendURL(ctx, &proto.ExtendURLRequest{
			Url:      shortenResp.Shortenurl,
			Username: "testuser",
		})
		if err != nil {
			log.Printf("ExtendURL failed: %v", err)
		} else {
			log.Printf("ExtendURL response: %v", extendResp)
		}
	}
}
