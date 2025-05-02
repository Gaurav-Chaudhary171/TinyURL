package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"TinyURL_Refactored/config"
	"TinyURL_Refactored/handlers"
	"TinyURL_Refactored/internal/service"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"

	"google.golang.org/grpc"
)

func main() {
	// Initialize database
	config.InitDB()

	// Auto migrate the schema
	err := config.DB.AutoMigrate(&model.Users{}, &model.GeneratedUrl{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create TinyURL service instance
	tinyURLService := service.NewTinyURLService()

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Printf("Failed to listen for gRPC: %v", err)
			return
		}

		s := grpc.NewServer()
		proto.RegisterTinyURLServiceServer(s, tinyURLService)

		log.Println("gRPC Server starting on :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("Failed to serve gRPC: %v", err)
		}
	}()

	// Setup HTTP handlers
	http.HandleFunc("/health", handlers.HealthCheck)
	http.HandleFunc("/v1/registeruser", handlers.RegisterHandler)
	http.HandleFunc("/v1/login", handlers.LoginHandler)
	http.HandleFunc("/v1/shortenurl", handlers.ShortenHandler)
	http.HandleFunc("/v1/extendurl", handlers.ExtendHandler)

	// Start HTTP server
	log.Println("HTTP Server starting on :8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Failed to serve HTTP: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down servers...")
}
