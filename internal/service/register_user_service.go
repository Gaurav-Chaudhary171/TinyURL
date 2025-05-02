// UserServer implements the gRPC user service
package service

import (
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"
	"context"
	"fmt"
)

type RegisterUserServer struct {
	proto.UnimplementedRegisterUserServer
}

// RegisterUser handles user registration
func (s *RegisterUserServer) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Check if username already exists
	var existingUser model.Users
	result := config.DB.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Create new user
	user := model.Users{
		Username: req.Username,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error registering user: %v", err)
	}

	return &proto.RegisterUserResponse{
		Status:   "success",
		Username: req.Username,
	}, nil
}

// generateUsername creates a username from first and last name
func generateUsername(firstName, lastName string) string {
	// Get first half of first name
	firstHalf := firstName
	if len(firstName) > 1 {
		firstHalf = firstName[:len(firstName)/2]
	}

	// Get last half of last name
	lastHalf := lastName
	if len(lastName) > 1 {
		lastHalf = lastName[len(lastName)/2:]
	}

	return firstHalf + lastHalf
}
