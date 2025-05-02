package service

import (
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"
	"context"
	"fmt"
)

type LoginServer struct {
	proto.UnimplementedTinyURLServiceServer
}

func (s *LoginServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user model.Users
	result := config.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found: %v", result.Error)
	}

	return &proto.LoginResponse{
		Status: "success",
		User: &proto.RegisterUserRequest{
			Username: user.Username,
		},
	}, nil
}
